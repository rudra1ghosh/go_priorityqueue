package workflows

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	pb "github.com/rudra1ghosh/HPE-Internship/project/requestmgmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.uber.org/cadence"
	"go.uber.org/cadence/activity"

	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

type RequestBody struct {
	WorkID string `json:"work_id"`
}

func init() {
	// Registering workflow and activtiy

	workflow.Register(CustomerWorkflow2)

	activity.Register(validateUser)

	activity.Register(Activity3)
	activity.Register(subscriptionDetails)
	activity.Register(blockStorage)

}

func CustomerWorkflow2(ctx workflow.Context, id int) error {

	retryPolicy := &cadence.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2,
		MaximumInterval:    time.Minute * 5,
		ExpirationInterval: time.Minute * 10,
		MaximumAttempts:    5,
	}
	currentState := "started Worklfow"
	err_c := workflow.SetQueryHandler(ctx, "current_state", func() (string, error) {
		return currentState, nil
	})
	if err_c != nil {
		currentState = "failed to register query handler"
		return err_c
	}

	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute * 60,
		StartToCloseTimeout:    time.Minute * 60,
		HeartbeatTimeout:       time.Minute * 60,
		RetryPolicy:            retryPolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// To get the Workflow ID and run ID
	wid := workflow.GetInfo(ctx).WorkflowExecution.ID
	rid := workflow.GetInfo(ctx).WorkflowExecution.RunID

	logger := workflow.GetLogger(ctx)

	logger.Info("Customer workflow started")
	var Result string

	err2 := workflow.ExecuteActivity(ctx, validateUser, wid).Get(ctx, &Result)
	currentState = "Validate User"
	if err2 != nil {
		logger.Error("Validate User Failed", zap.Error(err2))
		return err2
	}
	err3 := workflow.ExecuteActivity(ctx, subscriptionDetails, wid).Get(ctx, &Result)
	currentState = "Getting Subscription Details"
	if err3 != nil {
		logger.Error("Subscription Failed", zap.Error(err3))
		return err3
	}

	err := workflow.ExecuteActivity(ctx, Activity3, wid, rid, id).Get(ctx, &Result)
	currentState = "Enqueuing Activtiy"
	if err != nil {
		logger.Error("Activity failed.", zap.Error(err))
		return err
	}

	err1 := workflow.ExecuteActivity(ctx, wait, wid).Get(ctx, &Result)
	currentState = "Waiting for signal"
	if err1 != nil {
		logger.Error("Activity wait failed.", zap.Error(err1))
		return err1
	}

	err4 := workflow.ExecuteActivity(ctx, blockStorage, wid).Get(ctx, &Result)
	currentState = "Block Storage Granted"
	if err4 != nil {
		logger.Error("Activity failed.", zap.Error(err4))
		return err4
	}
	currentState = "Workflow completed"
	logger.Info("Workflow completed.", zap.String("Result", Result))

	return nil
}

func validateUser(ctx context.Context, workflow_id string) error {

	request1 := RequestBody{
		WorkID: workflow_id,
	}

	time.Sleep(time.Second * 5)

	requestBodyBytes, err := json.Marshal(request1)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("http://localhost:9090/validate", "application/json", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil

}

func blockStorage(ctx context.Context, workflow_id string) error {

	time.Sleep(time.Second * 5)
	endpoint := "http://localhost:9090/blockstorage"
	resp, err := http.Post(endpoint, "application/x-www-form-urlencoded", bytes.NewBufferString(""))
	if err != nil {
		fmt.Printf("Error posting to %s: %v\n", endpoint, err)
		return err

	}

	defer resp.Body.Close()

	return nil

}

func subscriptionDetails(ctx context.Context, workflow_id string) error {

	time.Sleep(time.Second * 5)
	endpoint := "http://localhost:9090/subscription"
	resp, err := http.Post(endpoint, "application/x-www-form-urlencoded", bytes.NewBufferString(""))
	if err != nil {
		fmt.Printf("Error posting to %s: %v\n", endpoint, err)
		return err

	}

	defer resp.Body.Close()

	return nil

}

func Activity3(ctx context.Context, workflow_id string, rid string, id int32) (string, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect %v", err)

	}

	defer conn.Close()
	c := pb.NewRequestManagementClient(conn)
	req := &pb.NewRequest{
		Wid: workflow_id, // workflow ID
		Rid: rid,         // run ID
		Id:  int32(id),   // request ID
	}

	ctx2, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel() // Ensure the context is canceled when done

	resp, err2 := c.CreateRequest(ctx2, req)
	if err2 != nil {
		log.Fatalf("failed to create request: %v", err2)
	}

	ans := fmt.Sprintf("Request created with Wid: %v", resp.GetWid())

	return ans, nil
}
