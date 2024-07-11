package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"sync"
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
	WorkID      string `json:"work_id"`
	PriorityID  int    `json:"p_id"`
	Flag        int
	DeTimestamp string `json:"timestamp"`
	EnTimestamp string
}

var dequeuedTasks []RequestBody
var mu sync.Mutex

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
}

var Q1 = Queue.Queue{
	Size: 100,
}

var Q2 = Queue.Queue{
	Size: 100,
}

var Q3 = Queue.Queue{
	Size: 100,
}

func (s *RequestManagementServer) CreateRequest(ctx context.Context, in *pb.NewRequest) (*pb.Request, error) {
	_, err := rpc1(ctx, in.GetWid(), in.GetRid(), in.GetId())
	if err != nil {
		log.Fatalf("error while enqueuing %v", err)
	}
	return &pb.Request{Wid: in.GetWid()}, nil
}

func rpc1(ctx context.Context, workflow_id string, rid string, id int32) (string, error) {
	switch id {
	case 1, 4:
		ans, err := rpc_function1(ctx, workflow_id, id, rid, &Q1)
		return ans, err
	case 2, 5:
		ans, err := rpc_function1(ctx, workflow_id, id, rid, &Q2)
		return ans, err
	case 3, 6:
		ans, err := rpc_function1(ctx, workflow_id, id, rid, &Q3)
		return ans, err
	default:
		return "Error Service not available", errors.New("unsupported service")
	}
}

func rpc_function1(ctx context.Context, workflow_id string, id int32, rid string, q *Queue.Queue) (string, error) {
	k := q.GetLength()

	switch id {
	case 1:
		customer1 := Queue.New(workflow_id, rid, ctx, id, time.Now(), k, 1)
		_, err := q.Enqueue(customer1)
		if err != nil {
			panic(err)
		}
	case 4:
		customer1 := Queue.New(workflow_id, rid, ctx, id, time.Now(), k, 0)
		_, err := q.Enqueue(customer1)
		if err != nil {
			panic(err)
		}
	case 2:
		customer1 := Queue.New(workflow_id, rid, ctx, id, time.Now(), k, 1)
		_, err := q.Enqueue(customer1)
		if err != nil {
			panic(err)
		}
	case 5:
		customer1 := Queue.New(workflow_id, rid, ctx, id, time.Now(), k, 0)
		_, err := q.Enqueue(customer1)
		if err != nil {
			panic(err)
		}
	case 3:
		customer1 := Queue.New(workflow_id, rid, ctx, id, time.Now(), k, 1)
		_, err := q.Enqueue(customer1)
		if err != nil {
			panic(err)
		}
	case 6:
		customer1 := Queue.New(workflow_id, rid, ctx, id, time.Now(), k, 0)
		_, err := q.Enqueue(customer1)
		if err != nil {
			panic(err)
		}
	}

	requestBody := RequestBody{
		WorkID:      workflow_id,
		PriorityID:  int(id),
		EnTimestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	mu.Lock()
	dequeuedTasks = append(dequeuedTasks, requestBody)
	mu.Unlock()

	return "", nil
}

func queueHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	queueContents := make([]Queue.Customer, len(Q1.Customers))
	copy(queueContents, Q1.Customers)

	data := struct {
		Customers []Queue.Customer
	}{
		Customers: queueContents,
	}

	funcMap := template.FuncMap{
		"isPrimePriority": isPP,
		"flagStatus":      flagStatus,
	}

	tmpl := template.Must(template.New("queue.html").Funcs(funcMap).ParseFiles("queue.html"))
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func queueHandler2(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	queueContents := make([]Queue.Customer, len(Q3.Customers))
	copy(queueContents, Q3.Customers)

	data := struct {
		Customers []Queue.Customer
	}{
		Customers: queueContents,
	}

	funcMap := template.FuncMap{
		"isPrimePriority": isPP,
		"flagStatus":      flagStatus,
	}

	tmpl := template.Must(template.New("queue.html").Funcs(funcMap).ParseFiles("queue.html"))
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func flagStatus(flag int) string {
	switch flag {
	case 0:
		return "Non-Prime"
	case 1:
		return "Prime"
	default:
		return "Unknown"
	}
}

func main() {
	go worker1()
	//go worker2()
	go worker3()

	go Q1.MoveToFrontIfOverdue(4)
	// go Q2.MoveToFrontIfOverdue(5)
	go Q3.MoveToFrontIfOverdue(6)

	log.Println("All workers ready")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterRequestManagementServer(s, &RequestManagementServer{})
	log.Printf("server listening at %v", lis.Addr())

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	http.HandleFunc("/dequeuedTasks", dequeuedTasksHandler)
	http.HandleFunc("/queue1", queueHandler)
	http.HandleFunc("/queue3", queueHandler2)
	http.ListenAndServe(":9092", nil)
	select {}
}

func isPP(priority int32) string {
	switch priority {
	case 1, 2, 3:
		return "Prime"
	case 4, 5, 6:
		return "Non-Prime"
	default:
		return "Unknown"
	}
}

func isPrimePriority(priorityID int) string {
	switch priorityID {
	case 1, 2, 3:
		return "Prime"
	case 4, 5, 6:
		return "Non-Prime"
	default:
		return "Unknown"
	}
}

func dequeuedTasksHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var networkingTasks, cloudEnterpriseTasks, blockStorageTasks []RequestBody
	for _, task := range dequeuedTasks {
		switch task.PriorityID {
		case 1, 4:
			networkingTasks = append(networkingTasks, task)
		case 2, 5:
			cloudEnterpriseTasks = append(cloudEnterpriseTasks, task)
		case 3, 6:
			blockStorageTasks = append(blockStorageTasks, task)
		}
	}

	data := struct {
		NetworkingTasks      []RequestBody
		CloudEnterpriseTasks []RequestBody
		BlockStorageTasks    []RequestBody
	}{
		NetworkingTasks:      networkingTasks,
		CloudEnterpriseTasks: cloudEnterpriseTasks,
		BlockStorageTasks:    blockStorageTasks,
	}

	tmpl := template.Must(template.New("dequeued_tasks.html").Funcs(template.FuncMap{
		"isPrimePriority": isPrimePriority,
		"flagStatus":      flagStatus,
	}).ParseFiles("dequeued_tasks.html"))

	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func worker1() {

	time.Sleep(1 * time.Minute)
	for c > -10 {
		if task_counter1 < 3 {
			task_counter1++
			go NetworkingServiceProcessing()
			time.Sleep(time.Millisecond * 100)
		} else {
			time.Sleep(time.Second)
		}
	}
}

func worker2() {
	for a > -10 {
		if task_counter2 < 3 {
			task_counter2++
			go PrivateCloudEnterpriseServiceProcessing()
			time.Sleep(time.Millisecond * 100)
		} else {
			time.Sleep(time.Second)
		}
	}
}

func worker3() {
	time.Sleep(1 * time.Minute)
	for b > -10 {
		if task_counter3 < 3 {
			task_counter3++
			go BlockStorageServiceProcessing()
			time.Sleep(time.Millisecond * 100)
		} else {
			time.Sleep(time.Second)
		}
	}
}

func NetworkingServiceProcessing() {
	response2, _, _, priority, _, err2, flag := Q1.Dequeue()
	if err2 == nil {
		fmt.Printf("DEQUEUED: WID: %s  Priority: %v Timestamp: %s\n", response2, priority, time.Now())
		mu.Lock()
		for i := range dequeuedTasks {
			if dequeuedTasks[i].WorkID == response2 {
				dequeuedTasks[i].DeTimestamp = time.Now().Format("2006-01-02 15:04:05")
				break
			}
		}
		mu.Unlock()
		request1 := RequestBody{
			WorkID:      response2,
			PriorityID:  int(priority),
			Flag:        flag,
			DeTimestamp: time.Now().Format("2006-01-02 15:04:05"),
		}

		time.Sleep(time.Second * 30)

		sendRequest(request1, "http://localhost:8090/endpoint1")

		time.Sleep(time.Millisecond)
	}
	task_counter1--
}

func PrivateCloudEnterpriseServiceProcessing() {
	response1, _, _, priority, _, err1, flag := Q2.Dequeue()
	if err1 == nil {
		fmt.Printf("DEQUEUED: WID: %s  Priority: %v Timestamp: %s\n", response1, priority, time.Now())

		mu.Lock()
		for i := range dequeuedTasks {
			if dequeuedTasks[i].WorkID == response1 {
				dequeuedTasks[i].DeTimestamp = time.Now().Format("2006-01-02 15:04:05")
				break
			}
		}
		mu.Unlock()
		request1 := RequestBody{
			WorkID:      response1,
			PriorityID:  int(priority),
			Flag:        flag,
			DeTimestamp: time.Now().Format("2006-01-02 15:04:05"),
		}

		time.Sleep(time.Second * 30)

		sendRequest(request1, "http://localhost:8090/endpoint2")

		time.Sleep(time.Millisecond)
	}
	task_counter2--
}

func BlockStorageServiceProcessing() {
	response3, _, _, priority, _, err3, flag := Q3.Dequeue()
	if err3 == nil {
		fmt.Printf("DEQUEUED: WID: %s  Priority: %v Timestamp: %s\n", response3, priority, time.Now())
		mu.Lock()
		for i := range dequeuedTasks {
			if dequeuedTasks[i].WorkID == response3 {
				dequeuedTasks[i].DeTimestamp = time.Now().Format("2006-01-02 15:04:05")
				break
			}
		}
		mu.Unlock()

		request1 := RequestBody{
			WorkID:      response3,
			PriorityID:  int(priority),
			Flag:        flag,
			DeTimestamp: time.Now().Format("2006-01-02 15:04:05"),
		}
		time.Sleep(time.Second * 5)
		sendRequest(request1, "http://localhost:8090/endpoint3")

		time.Sleep(time.Millisecond)
	}
	task_counter3--
}
