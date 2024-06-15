package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

const cadenceCLIImage = "ubercadence/cli:master"
const cadenceAddress = "host.docker.internal:7933"
const domain = "day56-domain"
const state = "current_state"

func main() {
	// Serve the index.html file from the current directory
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)
	http.HandleFunc("/submit", submitHandler)

	fmt.Println("Starting server at port 9095")
	err := http.ListenAndServe(":9095", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	workflowID := r.FormValue("workflow_id")
	runID := r.FormValue("run_id")

	dockerCmd := fmt.Sprintf("docker run --rm %s --address %s --domain %s workflow query -w %s -r %s -qt %s", cadenceCLIImage, cadenceAddress, domain, workflowID, runID, state)
	ans := executeCommand(dockerCmd)

	htmlResponse := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Workflow Status</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #f0f0f0;
					padding: 20px;
				}
				.container {
					max-width: 600px;
					margin: 0 auto;
					background-color: #fff;
					padding: 20px;
					border-radius: 8px;
					box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
				}
				h2 {
					color: #333;
					margin-bottom: 20px;
				}
				p {
					color: #666;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h2>Workflow Status</h2>
				<p><strong>Current Status:</strong> %s</p>
			</div>
		</body>
		</html>
	`

	fmt.Fprintf(w, htmlResponse, ans)
}

func executeCommand(command string) string {
	cmd := exec.Command("cmd", "/c", command)

	//cmd.Stderr = os.Stderr

	output, err := cmd.CombinedOutput()
	if err != nil {

		return "INVALID WORKLFOW AUR RUN ID"

	}

	return string(output)
}
