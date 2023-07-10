package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	tracing "github.com/sudo-NithishKarthik/syscalls-tracer/pkg/tracing"
)

func startTraceHandler(t *tracing.Tracer) func(w http.ResponseWriter, r *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
  // decode data from body
  decoder := json.NewDecoder(r.Body)
  var traceConf tracing.TracingConfiguration
  err := decoder.Decode(&traceConf)
  if err != nil {
      fmt.Printf("Unable to unmarshal the JSON config: %s", err)
      return
  }
  fmt.Println("Updating tracing configuration...")
  // set tracer config
  err = t.UpdateConfig(traceConf)
  if err != nil {
      fmt.Printf("Unable to create falco rule from the config: %s", err)
      return
  }
  fmt.Println("Starting to trace syscalls...")
  // start tracer
  err = t.Start()
  if err != nil {
      fmt.Printf("Unable to start tracer: %s", err)
      panic(err)
  }
  }
}

func stopTraceHandler(t *tracing.Tracer) func(w http.ResponseWriter, r *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Received request to stop tracing...")
	// Stop falco 
  err := t.Stop()
  if err != nil {
      panic(err)
  }
  }
}

func syscallsDataHandler(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "/falco/data.json")
}

func main() {
  tracer, err := tracing.NewTracer()
  if err != nil {
      panic(err)
  }
	http.HandleFunc("/start", startTraceHandler(&tracer))
	http.HandleFunc("/stop", stopTraceHandler(&tracer))
	http.HandleFunc("/data.json", syscallsDataHandler)
  fmt.Println("Starting server at locahost:9842...")
	log.Fatal(http.ListenAndServe(":9842", nil))
}
