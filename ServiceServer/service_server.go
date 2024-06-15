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
	WorkID     string `json:"work_id"`
	PriorityID int    `json:"p_id"`
}
type progress struct {
	wid string
	pid int
}

var progress1 []progress
var progress2 []progress
var progress3 []progress

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

func add(newWid string, newPid int, parray []progress) {
	prog := progress{
		wid: newWid,
		pid: newPid,
	}
	parray = append(parray, prog)
}

func progress1_fn(w http.ResponseWriter, r *http.Request) {

	for _, prog := range progress1 {
		fmt.Fprintf(w, "WORKFLOW ID: %s PRIORITY: %d\n", prog.wid, prog.pid)
	}

}
func progress2_fn(w http.ResponseWriter, r *http.Request) {

	for _, prog := range progress2 {
		fmt.Fprintf(w, "WORKFLOW ID: %s PRIORITY: %d\n", prog.wid, prog.pid)
	}

}

func progress3_fn(w http.ResponseWriter, r *http.Request) {

	for _, prog := range progress3 {
		fmt.Fprintf(w, "WORKFLOW ID: %s PRIORITY: %d\n", prog.wid, prog.pid)
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
	add(wid, runid, progress1)
	fmt.Println(wid)
	x := "PRIVATE CLOUD ENTERPRISE"

	fmt.Fprintf(w, "Service Name = %s\n", x)
	fmt.Fprintf(w, "Workflow ID= %s\n", wid)
	fmt.Fprintf(w, "Run ID = %d\n", runid)

}

func handleRequest2(w http.ResponseWriter, r *http.Request) {
	wid, runid := handleRequest(w, r, "Endpoint 2")
	add(wid, runid, progress2)
	x := "NETWORKING SERVICE"

	fmt.Fprintf(w, "Service Name = %s\n", x)
	fmt.Fprintf(w, "Workflow ID= %s\n", wid)
	fmt.Fprintf(w, "Run ID = %d\n", runid)
}

func handleRequest3(w http.ResponseWriter, r *http.Request) {
	wid, runid := handleRequest(w, r, "Endpoint 3")
	add(wid, runid, progress3)
	x := "BLOCK STORAGE SERVCIE"

	fmt.Fprintf(w, "Service Name = %s\n", x)
	fmt.Fprintf(w, "Workflow ID= %s\n", wid)
	fmt.Fprintf(w, "Run ID = %d\n", runid)
}

func handleRequest(w http.ResponseWriter, r *http.Request, endpoint string) (string, int) {

	var requestBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return "", 0
	}

	fmt.Printf("Received request at %s: WorkflowID=%s, Priority ID=%d\n", endpoint, requestBody.WorkID, requestBody.PriorityID)

	switch requestBody.PriorityID {

	case 1, 4:
		add(requestBody.WorkID, requestBody.PriorityID, progress1)

	case 2, 5:
		add(requestBody.WorkID, requestBody.PriorityID, progress2)

	case 3, 6:
		add(requestBody.WorkID, requestBody.PriorityID, progress3)
	}

	go triggerSignal(requestBody.WorkID)

	return requestBody.WorkID, requestBody.PriorityID

}

func main() {

	http.HandleFunc("/endpoint1", handleRequest1)
	http.HandleFunc("/endpoint2", handleRequest2)
	http.HandleFunc("/endpoint3", handleRequest3)
	http.HandleFunc("/progress1", progress1_fn)
	http.HandleFunc("/progress2", progress2_fn)
	http.HandleFunc("/progress3", progress3_fn)
	http.HandleFunc("/set-signal", setSignalHandler)
	http.HandleFunc("/get-signal", getSignalHandler)
	http.HandleFunc("/ws", wsHandler)

	fmt.Println("Server listening on port 8090...")
	log.Fatal(http.ListenAndServe(":8090", nil))

}
