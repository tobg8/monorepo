package main

import (
	_ "embed"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	sl "github.com/monorepo/domains/sample_library"
)

func YourHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request")
	log.Print(r)
	_, err := w.Write([]byte(sl.SampleGoFunction("Akhilesh")))
	if err != nil {
		return
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "8080"
	}

	return port
}

func main() {
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", YourHandler)
	// Bind to a port and pass our router in
	port := getPort()
	log.Println("Going to listen on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}
