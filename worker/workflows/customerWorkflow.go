package workflows

import (
	"bytes"
	"context"
	"fmt"

	"log"
	"net/http"

	"time"

	"github.com/gorilla/websocket"

	pb "github.com/rudra1ghosh/HPE-Internship/project/requestmgmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.uber.org/cadence"
	"go.uber.org/cadence/activity"

	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

func init() {
	// Registering workflow and activtiy

	workflow.Register(CustomerWorkflow)
	activity.Register(Activity1)
	activity.Register(Application_Details)
	activity.Register(Quiesce)
	activity.Register(snapshot)
	activity.Register(UnQuiesce)
	activity.Register(backup)
	activity.Register(wait)

}

const (
	address = "localhost:50051"
)

func CustomerWorkflow(ctx workflow.Context, id int) error {

	retryPolicy := &cadence.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2,
		MaximumInterval:    time.Minute * 5,
		ExpirationInterval: time.Minute * 10,
		MaximumAttempts:    5,
	}

	currentState := "started Worklfow"
	err := workflow.SetQueryHandler(ctx, "current_state", func() (string, error) {
		return currentState, nil
	})
	if err != nil {
		currentState = "failed to register query handler"
		return err
	}

	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute * 5,
		StartToCloseTimeout:    time.Minute * 5,
		HeartbeatTimeout:       time.Minute * 5,
		RetryPolicy:            retryPolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	currentState = "waiting to be started"
	// To get the Workflow ID and run ID
	wid := workflow.GetInfo(ctx).WorkflowExecution.ID
	rid := workflow.GetInfo(ctx).WorkflowExecution.RunID

	logger := workflow.GetLogger(ctx)
	logger.Info("Customer workflow started")
	var Result string

	if id == 1 || id == 4 {

		err1 := workflow.ExecuteActivity(ctx, Application_Details, wid).Get(ctx, &Result)
		currentState = "Application Details"
		if err1 != nil {
			logger.Error("Application Failed", zap.Error(err1))
			currentState = "Application Details Error"
			return err1
		}

		err2 := workflow.ExecuteActivity(ctx, Quiesce, wid).Get(ctx, &Result)
		currentState = "Quiesce Details"
		if err2 != nil {
			logger.Error("Quiece Failed", zap.Error(err2))
			return err2
		}

		err := workflow.ExecuteActivity(ctx, Activity1, wid, rid, id).Get(ctx, &Result)

		currentState = "Request enqueued "
		if err != nil {
			logger.Error("Activity Enqueue failed.", zap.Error(err))
			return err
		}

		err4 := workflow.ExecuteActivity(ctx, wait, wid).Get(ctx, &Result)
		currentState = "waiting for signal"
		if err4 != nil {
			logger.Error("Activity wait failed.", zap.Error(err4))
			return err4
		}
		err3 := workflow.ExecuteActivity(ctx, UnQuiesce, wid).Get(ctx, &Result)
		currentState = "unquiece activity "
		if err3 != nil {
			logger.Error("Unquiece Failed", zap.Error(err3))
			return err3
		}
		logger.Info("Workflow completed.", zap.String("Result", Result))
		currentState = "Completed "
		return nil
	}

	if id == 2 || id == 5 {

		err1 := workflow.ExecuteActivity(ctx, Application_Details, wid).Get(ctx, &Result)
		currentState = "Application Details"
		if err1 != nil {
			logger.Error("Application Failed", zap.Error(err1))
			return err1
		}

		err2 := workflow.ExecuteActivity(ctx, snapshot, wid).Get(ctx, &Result)
		currentState = "Taking a Snapshot"
		if err2 != nil {
			logger.Error("Quiece Failed", zap.Error(err2))
			return err2
		}

		err := workflow.ExecuteActivity(ctx, Activity1, wid, rid, id).Get(ctx, &Result)
		currentState = "Enqueuing the snapshot Activity"
		if err != nil {
			logger.Error("Activity Enqueue failed.", zap.Error(err))
			return err
		}

		err4 := workflow.ExecuteActivity(ctx, wait, wid).Get(ctx, &Result)
		currentState = "Waiting for signal"
		if err4 != nil {
			logger.Error("Activity wait failed.", zap.Error(err4))
			return err4
		}

		err3 := workflow.ExecuteActivity(ctx, backup, wid).Get(ctx, &Result)
		currentState = "Creating backup"
		if err3 != nil {
			logger.Error("Subscription Failed", zap.Error(err3))
			return err3
		}
		currentState = "Workflow completetd "
		logger.Info("Workflow completed.", zap.String("Result", Result))

		return nil

	}

	logger.Info("Workflow completed.", zap.String("Result", Result))

	return nil
}

func Application_Details(ctx context.Context, workflow_id string) error {

	time.Sleep(time.Second * 5)
	endpoint := "http://localhost:9090/details"
	resp, err := http.Post(endpoint, "application/x-www-form-urlencoded", bytes.NewBufferString(""))
	if err != nil {
		fmt.Printf("Error posting to %s: %v\n", endpoint, err)
		return err

	}

	defer resp.Body.Close()

	return nil

}

func Quiesce(ctx context.Context, workflow_id string) error {

	time.Sleep(time.Second * 5)
	endpoint := "http://localhost:9090/Quiesce"
	resp, err := http.Post(endpoint, "application/x-www-form-urlencoded", bytes.NewBufferString(""))
	if err != nil {
		fmt.Printf("Error posting to %s: %v\n", endpoint, err)
		return err

	}

	defer resp.Body.Close()

	return nil

}

func snapshot(ctx context.Context, workflow_id string) error {

	time.Sleep(time.Second * 5)
	endpoint := "http://localhost:9090/snapshot"
	resp, err := http.Post(endpoint, "application/x-www-form-urlencoded", bytes.NewBufferString(""))
	if err != nil {
		fmt.Printf("Error posting to %s: %v\n", endpoint, err)
		return err

	}

	defer resp.Body.Close()

	return nil

}

func UnQuiesce(ctx context.Context, workflow_id string) error {

	time.Sleep(time.Second * 5)
	endpoint := "http://localhost:9090/UnQuiesce"
	resp, err := http.Post(endpoint, "application/x-www-form-urlencoded", bytes.NewBufferString(""))
	if err != nil {
		fmt.Printf("Error posting to %s: %v\n", endpoint, err)
		return err

	}

	defer resp.Body.Close()

	return nil

}

func backup(ctx context.Context, workflow_id string) error {

	time.Sleep(time.Second * 5)
	endpoint := "http://localhost:9090/backup"
	resp, err := http.Post(endpoint, "application/x-www-form-urlencoded", bytes.NewBufferString(""))
	if err != nil {
		fmt.Printf("Error posting to %s: %v\n", endpoint, err)
		return err

	}

	defer resp.Body.Close()

	return nil

}

func wait(ctx context.Context, workflow_id string) error {

	var expectedSignal = workflow_id
	waitingFunction(expectedSignal)

	return nil
}

func Activity1(ctx context.Context, workflow_id string, rid string, id int32) (string, error) {
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

func waitingFunction(expectedSignal string) {
	fmt.Println("Waiting for the signal...")

	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8090/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	for {

		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		if string(message) == expectedSignal {
			fmt.Println("Signal received, continuing execution...")
			return
		}
	}
}
