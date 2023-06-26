package main

import (
	"fmt"
	"log"
	"net/http"
)

func startTraceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello")
	// Run falco with appropriate configuration
}

func stopTracingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello")
	// Stop falco and send the resultant syscalls
}

func main() {
	http.HandleFunc("/start", startTraceHandler)
	http.HandleFunc("/stop", stopTracingHandler)
	log.Fatal(http.ListenAndServe(":9842", nil))
}
