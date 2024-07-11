package main

import (
	"fmt"
	"html/template"
	"strconv"

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
		fmt.Fprintf(w, "Parse form error: %v", err)
		return
	}

	service_request := r.FormValue("serviceId")
	num, _ := strconv.Atoi(service_request)

	htmlTemplate := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Workflow Information</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f9;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
        }
        .container {
            background-color: #f0f8f9;
            border-radius: 8px;
            box-shadow: 0 0 20px rgba(0, 0, 0, 0.1);
            width: 500px;
            padding: 20px;
            max-width: 100%;
            overflow: hidden;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-bottom: 20px;
            table-layout: fixed;
        }
        th, td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #ddd;
            word-wrap: break-word;
        }
        th {
            background-color: #7ac142;
            color: white;
            white-space: nowrap;
        }
        .highlight {
            font-weight: bold;
            color: #2c3e50;
        }
        .track-button {
            display: block;
            width: 100%;
            text-align: center;
            margin-top: 20px;
        }
        .track-button button {
            background-color: #7ac142;
            color: white;
            border: none;
            padding: 12px 20px;
            border-radius: 4px;
            font-size: 16px;
            cursor: pointer;
            transition: background-color 0.3s ease;
        }
        .track-button button:hover {
            background-color: #68a12d;
        }
    </style>
</head>
<body>
    <div class="container">
        <h2 style="text-align: center; color: #7ac142; margin-bottom: 20px;">Workflow Information</h2>
        <table>
            <tr>
                <th style="min-width: 150px;">Parameter</th>
                <th>Value</th>
            </tr>
            <tr>
                <td style="word-break: break-all;">Workflow ID:</td>
                <td><span class="highlight">{{.WorkflowID}}</span></td>
            </tr>
            <tr>
                <td style="word-break: break-all;">Run ID:</td>
                <td><span class="highlight">{{.RunID}}</span></td>
            </tr>
        </table>
        <form action="http://localhost:9095/submit" method="POST" class="track-button">
            <input type="hidden" name="workflow_id" value="{{.WorkflowID}}">
            <input type="hidden" name="run_id" value="{{.RunID}}">
            <button type="submit">Track</button>
        </form>
    </div>
</body>
</html>
	`

	var wid, rid string
	switch num {
	case 1:
		wid, rid = h.start_worklfow(1)
	case 2:
		wid, rid = h.start_worklfow(2)
	case 3:
		wid, rid = h.start_worklfow2(3)
	case 4:
		wid, rid = h.start_worklfow(4)
	case 5:
		wid, rid = h.start_worklfow(5)
	case 6:
		wid, rid = h.start_worklfow2(6)
	}

	data := struct {
		WorkflowID string
		RunID      string
	}{
		WorkflowID: wid,
		RunID:      rid,
	}

	tmpl := template.Must(template.New("index").Parse(htmlTemplate))

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Service) start_worklfow(id int) (string, string) {

	wo := client.StartWorkflowOptions{
		TaskList:                     taskList,
		ExecutionStartToCloseTimeout: time.Minute * 30,
	}

	req, err := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.CustomerWorkflow, 1)

	if err != nil {

		h.logger.Error("Service not available ")
		return "", ""

	}

	_, err2 := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.CustomerWorkflow, 1)

	if err2 != nil {

		h.logger.Error("Service not available ")
		return "", ""
	}
	_, err3 := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.CustomerWorkflow, 1)

	if err3 != nil {

		h.logger.Error("Service not available ")
		return "", ""

	}

	_, err4 := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.CustomerWorkflow, 4)

	if err4 != nil {

		h.logger.Error("Service not available ")
		return "", ""

	}
	_, err5 := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.CustomerWorkflow, 1)

	if err5 != nil {

		h.logger.Error("Service not available ")
		return "", ""

	}
	_, err6 := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.CustomerWorkflow, 1)

	if err6 != nil {

		h.logger.Error("Service not available ")
		return "", ""

	}

	_, err7 := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.CustomerWorkflow, 4)

	if err7 != nil {

		h.logger.Error("Service not available ")
		return "", ""

	}
	_, err8 := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.CustomerWorkflow, 4)

	if err8 != nil {

		h.logger.Error("Service not available ")
		return "", ""

	}
	_, err9 := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.CustomerWorkflow, 1)

	if err9 != nil {

		h.logger.Error("Service not available ")
		return "", ""

	}
	_, err10 := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.CustomerWorkflow, 4)

	if err10 != nil {

		h.logger.Error("Service not available ")
		return "", ""

	}

	return req.ID, req.RunID

}

func (h *Service) start_worklfow2(id int) (string, string) {

	wo := client.StartWorkflowOptions{
		TaskList:                     taskList2,
		ExecutionStartToCloseTimeout: time.Minute * 50,
	}

	ans, err := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.CustomerWorkflow2, 6)

	if err != nil {

		return "", ""

	}
	_, err9 := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.CustomerWorkflow2, 6)

	if err9 != nil {

		h.logger.Error("Service not available ")
		return "", ""

	}
	_, err10 := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.CustomerWorkflow2, 3)

	if err10 != nil {

		h.logger.Error("Service not available ")
		return "", ""

	}
	_, err13 := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.CustomerWorkflow2, 3)

	if err13 != nil {

		h.logger.Error("Service not available ")
		return "", ""

	}
	_, err14 := h.cadenceAdapter.CadenceClient.StartWorkflow(context.Background(), wo, workflows.CustomerWorkflow2, 3)

	if err14 != nil {

		h.logger.Error("Service not available ")
		return "", ""

	}
	return ans.ID, ans.RunID

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
