package main

import (
	"bytes"
	"context"
	"log"

	"encoding/json"
	"errors"
	"fmt"

	"net"
	"net/http"
	"time"

	"github.com/rudra1ghosh/HPE-Internship/project/Queue"
	pb "github.com/rudra1ghosh/HPE-Internship/project/requestmgmt"
	"google.golang.org/grpc"
)

var task_counter1 = 0
var task_counter2 = 0
var task_counter3 = 0
var c = 0
var a = 0
var b = 0

const (
	port = ":50051"
)

type RequestManagementServer struct {
	pb.UnimplementedRequestManagementServer
}

type RequestBody struct {
	WorkID string `json:"work_id"`
	RunID  string `json:"run_id"`
}

func sendRequest(requestBody RequestBody, url string) {
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	//	fmt.Println(requestBody.WorkID)
	//	fmt.Println(requestBody.RunID)

	//fmt.Println("Response from server:")
	//fmt.Println(resp.Status)
}

var Q1 = Queue.Queue{

	Size: 30,
}

var Q2 = Queue.Queue{

	Size: 30,
}
var Q3 = Queue.Queue{

	Size: 30,
}

func (s *RequestManagementServer) CreateRequest(ctx context.Context, in *pb.NewRequest) (*pb.Request, error) {

	//log.Printf("received: %v", in.GetWid())
	_, err := rpc1(ctx, in.GetWid(), in.GetRid(), in.GetId())

	if err != nil {
		log.Fatalf("error while enquing %v", err)
	}
	return &pb.Request{Wid: in.GetWid()}, nil

}
func rpc1(ctx context.Context, workflow_id string, rid string, id int32) (string, error) {
	switch id {
	case 1, 4:
		ans, err := rpc_function1(ctx, workflow_id, id, rid, &Q1)
		Q1.SortCustomers()

		return ans, err
	case 2, 5:
		ans, err := rpc_function1(ctx, workflow_id, id, rid, &Q2)
		Q2.SortCustomers()

		return ans, err
	case 3, 6:
		ans, err := rpc_function1(ctx, workflow_id, id, rid, &Q3)
		Q3.SortCustomers()

		return ans, err
	default:
		return "Error Service not available ", errors.ErrUnsupported
	}

}
func rpc_function1(ctx context.Context, workflow_id string, id int32, rid string, q *Queue.Queue) (string, error) {

	for q.GetLength() >= q.Size/2 {
		time.Sleep(time.Second)
	}

	customer1 := Queue.New(workflow_id, rid, ctx, id)
	ans := fmt.Sprintf("Enqueud at %s", time.Now())

	//fmt.Printf("WID: %s  Priority: %v\n", workflow_id, id)
	_, err := q.Enqueue(customer1)
	if err != nil {
		panic(err)
	}
	return ans, err
}

func main() {

	go worker1()
	go worker2()
	go worker3()

	log.Println("All workers ready ")
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen : %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterRequestManagementServer(s, &RequestManagementServer{})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve at: %v", err)
	}

	select {}
}

func worker1() {

	for c > -1 {
		if task_counter1 < 15 {
			task_counter1++
			go PrivateCloudEnterpriseServiceProcessing()
			time.Sleep(time.Millisecond * 100)
		} else {
			time.Sleep(time.Second)
		}

	}
}
func worker2() {

	for a > -1 {
		if task_counter2 < 15 {

			task_counter2++
			go NetworkingServiceProcessing()
			time.Sleep(time.Millisecond * 100)
		} else {
			time.Sleep(time.Second)
		}

	}

}

func worker3() {

	for b > -1 {
		if task_counter3 < 15 {
			task_counter3++
			go BlockStorageServiceProcessing()
			time.Sleep(time.Millisecond * 100)
		} else {
			time.Sleep(time.Second)
		}

	}

}

func PrivateCloudEnterpriseServiceProcessing() {

	response1, runid1, _, priority, err1 := Q1.Dequeue()

	if err1 == nil {
		fmt.Printf("WID: %s  Priority: %v\n", response1, priority)
		// time waiting for completion of service
		time.Sleep(time.Second * 10)

		request1 := RequestBody{
			WorkID: response1,
			RunID:  runid1,
		}
		sendRequest(request1, "http://localhost:8090/endpoint1")

		time.Sleep(time.Millisecond)

	}
	task_counter1--
}
func NetworkingServiceProcessing() {

	response2, runid2, _, priority, err2 := Q2.Dequeue()

	if err2 == nil {
		fmt.Printf("WID: %s  Priority: %v\n", response2, priority)
		time.Sleep(time.Second * 10)

		request1 := RequestBody{
			WorkID: response2,
			RunID:  runid2,
		}
		sendRequest(request1, "http://localhost:8090/endpoint2")
		time.Sleep(time.Millisecond)

	}
	task_counter2--

}

func BlockStorageServiceProcessing() {

	response3, runid3, _, priority, err3 := Q3.Dequeue()

	if err3 == nil {
		fmt.Printf("WID: %s  Priority: %v\n", response3, priority)
		time.Sleep(time.Second * 10)

		request1 := RequestBody{
			WorkID: response3,
			RunID:  runid3,
		}
		sendRequest(request1, "http://localhost:8090/endpoint3")

		time.Sleep(time.Millisecond)

	}
	task_counter3--
}
