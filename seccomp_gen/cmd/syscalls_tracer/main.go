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
      e := fmt.Sprintf( "Unable to unmarshal the JSON config: %s", err)
      http.Error(w, e, 400)
			return
		}
		fmt.Println("Updating tracing configuration...")
		// set tracer config
		err = t.UpdateConfig(traceConf)
		if err != nil {
      e := fmt.Sprintf("Unable to create falco rule from the config: %s", err)
      http.Error(w, e, 500)
			return
		}
    // verify that there are no running traces
    /**
     * We can have multiple falco processes running simultaneously tracing 
     * multiple different entities. Therefore, we should be able to support tracing
     * multiple entities at the same time. But we will keep things simple for now a
     * and don't allow multiple simultaneous traces.
    **/
    if t.IsTracing {
      e := fmt.Sprintf("There is a trace running already. Multiple simultaneous tracing is not yet supported.")
      http.Error(w, e, 500)
      return
    }
		fmt.Println("Starting to trace syscalls...")
		// start tracer
		err = t.Start()
		if err != nil {
      e := fmt.Sprintf("Unable to start tracer: %s", err)
      http.Error(w, e, 500)
      return
		}
    t.IsTracing = true
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
    t.IsTracing = false
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
