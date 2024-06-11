package main

import (
	"fmt"
	"net/http"

	"github.com/brianvoe/gofakeit/v6"
)

func applicationDetails(w http.ResponseWriter, r *http.Request) {

	gofakeit.Seed(0)

	// Generate a fake name and phone number
	name := gofakeit.Name()
	phone := gofakeit.Phone()

	fmt.Fprintf(w, "NAME: %s\n", name)
	fmt.Fprintf(w, "PHONE: %s", phone)

}

func validateUser(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "User has been Verified")

}

func getSubscription(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Getting Subscription Details: PRIME || NON-PRIME")

}
func Quiesce(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Quiesce Process started")

}

func storageSnapshot(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Creating storage snapshot")

}

func unQuiesce(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "UnQuiesce Process started")
}

func backup(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Catalogue Backup")

}

func blockStorage(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Block Storage granted")

}

func main() {
	// Workflow 1

	http.HandleFunc("/details", applicationDetails)
	http.HandleFunc("/Quiesce", Quiesce)
	http.HandleFunc("/snapshot", storageSnapshot)
	http.HandleFunc("/UnQuiesce", unQuiesce)
	http.HandleFunc("/backup", backup)

	// Workflow 2

	http.HandleFunc("/subscription", getSubscription)
	http.HandleFunc("/validate", validateUser)
	http.HandleFunc("/blockstorage", blockStorage)

	fmt.Println("Starting server on :9090...")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}

}
