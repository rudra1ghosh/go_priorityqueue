package main

import (
	"fmt"

	"context"
	"log"
	"net/http"

	"time"

	"github.com/rudra1ghosh/HPE-Internship/project/adapters/cadenceAdapter"
	"github.com/rudra1ghosh/HPE-Internship/project/config"

	"github.com/rudra1ghosh/HPE-Internship/project/worker/workflows"
	"go.uber.org/cadence/client"
	"go.uber.org/cadence/worker"
	"go.uber.org/zap"
)

const cadenceCLIImage = "ubercadence/cli:master"
const cadenceAddress = "host.docker.internal:7933"
const domain = "day56-domain"
const taskList = "Service_process"
const taskList2 = "Service2_process"

const (
	address = "localhost:50051"
)

type Service struct {
	cadenceAdapter *cadenceAdapter.CadenceAdapter
	logger         *zap.Logger
}

func (h *Service) formHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "Parse form error : %v", err)
		return
	}

	service_request := r.FormValue("service")
	fmt.Fprintf(w, "Request For %s Service Submitted\n", service_request)

	switch service_request {

	case "networking_prime":

		ans := h.start_worklfow(1)
		if ans == false {
			return
		}

	case "networking_non_prime":
		ans := h.start_worklfow(4)
		if ans == false {
			return
		}
	case "cloud_prime":
		ans := h.start_worklfow(2)
		if ans == false {
			return
		}
	case "cloud_non_prime":
		ans := h.start_worklfow(5)
		if ans == false {
			return
		}

	case "storage_prime":
		ans := h.start_worklfow2(3)
		if ans == false {
			return
		}

	case "storage_non_prime":
		ans := h.start_worklfow2(6)
		if ans == false {
			return
		}

	}

}

func (h *Service) start_worklfow(id int) bool {

	wo := client.StartWorkflowOptions{
		TaskList:                     taskList,
		ExecutionStartToCloseTimeout: time.Hour * 24,
	}

	workflownum, err := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.CustomerWorkflow, id)

	fmt.Printf("WORKFLOW ID: %s PRIORITY: %v\n", workflownum.ID, id)

	if err != nil {

		h.logger.Error("Service not available ")
		return false

	}

	return true

}
func (h *Service) start_worklfow2(id int) bool {

	wo := client.StartWorkflowOptions{
		TaskList:                     taskList2,
		ExecutionStartToCloseTimeout: time.Hour * 24,
	}

	workflownum, err := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.CustomerWorkflow2, id)

	fmt.Printf("WORKFLOW ID: %s PRIORITY: %v\n", workflownum.ID, id)

	if err != nil {

		h.logger.Error("Service not available ")
		return false

	}

	return true

}

func startWorkers(h *cadenceAdapter.CadenceAdapter, taskList string) {
	// Configure worker options.
	workerOptions := worker.Options{
		MetricsScope: h.Scope,
		Logger:       h.Logger,
	}

	cadenceWorker := worker.New(h.ServiceClient, h.Config.Domain, taskList, workerOptions)
	err := cadenceWorker.Start()
	if err != nil {
		h.Logger.Error("Failed to start workers.", zap.Error(err))
		panic("Failed to start workers")

	}

}

var appConfig config.AppConfig
var cadenceClient cadenceAdapter.CadenceAdapter
var service = Service{&cadenceClient, appConfig.Logger}

func main() {

	fmt.Println("Starting Worker..")

	appConfig.Setup()

	cadenceClient.Setup(&appConfig.Cadence)

	startWorkers(&cadenceClient, taskList)
	startWorkers(&cadenceClient, taskList2)

	fmt.Println("Cadence worker ready ")

	fileServer := http.FileServer(http.Dir("./static"))

	http.Handle("/", fileServer)

	http.HandleFunc("/form", service.formHandler)

	fmt.Printf("Starting server at 8080 port\n")

	err2 := http.ListenAndServe(":8080", nil)

	if err2 != nil {
		log.Fatal(err2)
	}

}
