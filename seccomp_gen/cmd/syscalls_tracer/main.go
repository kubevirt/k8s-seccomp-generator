package main

import (
	"encoding/json"
	"log"
	"net/http"

	tracing "github.com/sudo-NithishKarthik/syscalls-tracer/pkg/tracing"
)

func startTraceHandler(t tracing.Tracer) func(w http.ResponseWriter, r *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
  // decode data from body
  decoder := json.NewDecoder(r.Body)
  var traceConf tracing.TracingConfiguration
  err := decoder.Decode(&traceConf)
  if err != nil {
      panic(err)
  }
  // set tracer config
  err = t.SetConfig(traceConf)
  if err != nil {
      panic(err)
  }
  // start tracer
  err = t.Start()
  if err != nil {
      panic(err)
  }
  }
}

func stopTraceHandler(t tracing.Tracer) func(w http.ResponseWriter, r *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
	// Stop falco and send the resultant syscalls data.json file
  err := t.Stop()
  if err != nil {
      panic(err)
  }
  }
}

func main() {
  tracer, err := tracing.NewTracer()
  if err != nil {
      panic(err)
  }
	http.HandleFunc("/start", startTraceHandler(tracer))
	http.HandleFunc("/stop", stopTraceHandler(tracer))
	log.Fatal(http.ListenAndServe(":9842", nil))
}
