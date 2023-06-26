package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func startTraceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello")
	// Run falco in such a way that it will start tracing with:
	// 1. The given set of rules
	// 2. Using program_output to stream the output to the falco-syscalls-formatter binary with appropriate arguments
	// 3. Kubernetes option enabled
	// 4. Set to json_output
	// With the following flags:
	// 1. -A for monitoring all the syscall events
	// 2. --k8s-api
	// 3. --k8s-api-cert
	// 4. --unbuffered
  os.Setenv("FALCO_BPF_PROBE", "/falco/falco-bpf.o")
  falcoCommand := exec.Command("/usr/bin/falco",
    "-A",
    "-U",
    "-r", "/falco/rules.yaml", 
    "-k", "https://$KUBERNETES_SERVICE_HOST",
    "-K", "/var/run/secrets/kubernetes.io/serviceaccount/token",
    "--option", "program_output.enabled=true",
    "--option", "program_output.keep_alive=true",
    "--option", "program_output.program=/falco/falco-syscalls-formatter",
    "--option", "stdout_output.enabled=false",
    "--option", "syslog_output.enabled=false",
    "--option", "file_output.enabled=false",
    "--option", "json_output=true",
    )

}

func stopTracingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello")
	// Stop falco and send the resultant syscalls data.json file
}

func main() {
	http.HandleFunc("/start", startTraceHandler)
	http.HandleFunc("/stop", stopTracingHandler)
	log.Fatal(http.ListenAndServe(":9842", nil))
}
