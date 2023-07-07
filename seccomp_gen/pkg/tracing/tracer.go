package tracing

import (
	"fmt"
	"os"
	"os/exec"
)

type TracingConfiguration struct {
  PodName string `json:"podName"`
  ContainerName string `json:"containerName"`
  PodLabel map[string]string `json:"podLabel"`
}

type Tracer struct { 
  falcoProcess *os.Process
  Config TracingConfiguration
}

func NewTracer() (Tracer, error) {
  tracer := Tracer{}
  return tracer, nil
}


func (t *Tracer) SetConfig(config TracingConfiguration) error {
  t.Config = config
  rule, err := CreateFalcoRule(config)
  if err != nil {
      return err
  }
  // Write the rule to /falco/rules.yaml
  err = os.WriteFile("/falco/rules.yaml", rule, 0644)
  if err != nil {
    return err
  }
  return nil
}

// Start Falco process and update the struct with the falco process
func (t *Tracer) Start() error {
  khost := os.Getenv("KUBERNETES_SERVICE_HOST")
  if khost == "" {
    return fmt.Errorf("KUBERNETES_SERVICE_HOST not set")
  }
  falcoCommand := exec.Command("/usr/bin/falco",
    "-A",
    "-U",
    "-r", "/falco/rules.yaml", 
    "-k", "https://"+khost,
    "-K", "/var/run/secrets/kubernetes.io/serviceaccount/token",
    // TODO: What if crio is not used?
    "--cri", "/var/run/crio/crio.sock", 
    "--option", "program_output.enabled=true",
    "--option", "program_output.keep_alive=true",
    "--option", "program_output.program=/falco/formatter",
    "--option", "stdout_output.enabled=false",
    "--option", "syslog_output.enabled=false",
    "--option", "file_output.enabled=true",
    "--option", "file_output.filename=/falco/logs.txt",
    "--option", "json_output=true",
    )
  falcoCommand.Env = os.Environ()
  falcoCommand.Env = append(falcoCommand.Env, "FALCO_BPF_PROBE=/falco/falco-bpf.o")
  f, _ := os.Create("/falco/cmd_logs.txt")
  falcoCommand.Stdout = f
  falcoCommand.Stderr = f
  // we have to call Process.Release when stopping it
  err := falcoCommand.Start()
  if err != nil {
      return err
  }
  t.falcoProcess = falcoCommand.Process
  fmt.Println("Falco process struct value: ", t.falcoProcess, " and the value returned by c.Process: ", falcoCommand.Process)
  return nil
}

// Stop the tracer by sending interrupt to the Falco process
func (t *Tracer) Stop() error {
  if t.falcoProcess == nil {
    return fmt.Errorf("Process not found in the struct")
  }
  err := (*t.falcoProcess).Signal(os.Interrupt)
  if err != nil {
    return err
  }
  t.falcoProcess.Release()
  t.falcoProcess = nil
  return nil
}
