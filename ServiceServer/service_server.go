package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const cadenceCLIImage = "ubercadence/cli:master"
const cadenceAddress = "host.docker.internal:7933"
const domain = "day56-domain"

type RequestBody struct {
	WorkID string `json:"work_id"`
	RunID  string `json:"run_id"`
}

var (
	signal    string
	mu        sync.Mutex
	clients   = make(map[*websocket.Conn]struct{})
	clientsMu sync.Mutex
	upgrader  = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func setSignalHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	newSignal := r.FormValue("signal")
	if newSignal == "" {
		http.Error(w, "Signal not provided", http.StatusBadRequest)
		return
	}

	mu.Lock()
	signal = newSignal
	mu.Unlock()

	notifyClients(newSignal)

	fmt.Fprintln(w, "Signal set to", newSignal)
}

func getSignalHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	currentSignal := signal
	mu.Unlock()

	fmt.Fprintln(w, currentSignal)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	clientsMu.Lock()
	clients[conn] = struct{}{}
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			return
		}
	}
}

func notifyClients(signal string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(signal))
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}

func triggerSignal(wid string) {

	fmt.Println("Triggering the signal...")
	resp, err := http.PostForm("http://localhost:8090/set-signal",
		map[string][]string{"signal": {wid}})
	if err != nil {
		fmt.Println("Error triggering signal:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Signal triggered successfully.")
}

func handleRequest1(w http.ResponseWriter, r *http.Request) {
	wid, runid := handleRequest(w, r, "Endpoint 1")
	fmt.Println(wid)
	x := "PRIVATE CLOUD ENTERPRISE"

	fmt.Fprintf(w, "Service Name = %s\n", x)
	fmt.Fprintf(w, "Workflow ID= %s\n", wid)
	fmt.Fprintf(w, "Run ID = %s\n", runid)

}

func handleRequest2(w http.ResponseWriter, r *http.Request) {
	wid, runid := handleRequest(w, r, "Endpoint 2")
	x := "NETWORKING SERVICE"

	fmt.Fprintf(w, "Service Name = %s\n", x)
	fmt.Fprintf(w, "Workflow ID= %s\n", wid)
	fmt.Fprintf(w, "Run ID = %s\n", runid)
}

func handleRequest3(w http.ResponseWriter, r *http.Request) {
	wid, runid := handleRequest(w, r, "Endpoint 3")
	x := "BLOCK STORAGE SERVCIE"

	fmt.Fprintf(w, "Service Name = %s\n", x)
	fmt.Fprintf(w, "Workflow ID= %s\n", wid)
	fmt.Fprintf(w, "Run ID = %s\n", runid)
}

func handleRequest(w http.ResponseWriter, r *http.Request, endpoint string) (string, string) {

	var requestBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return "", ""
	}

	fmt.Printf("Received request at %s: WorkflowID=%s, RunID=%s\n", endpoint, requestBody.WorkID, requestBody.RunID)

	go triggerSignal(requestBody.WorkID)

	return requestBody.WorkID, requestBody.RunID

}

func main() {

	http.HandleFunc("/endpoint1", handleRequest1)
	http.HandleFunc("/endpoint2", handleRequest2)
	http.HandleFunc("/endpoint3", handleRequest3)

	http.HandleFunc("/set-signal", setSignalHandler)
	http.HandleFunc("/get-signal", getSignalHandler)
	http.HandleFunc("/ws", wsHandler)

	fmt.Println("Server listening on port 8090...")
	log.Fatal(http.ListenAndServe(":8090", nil))

}
