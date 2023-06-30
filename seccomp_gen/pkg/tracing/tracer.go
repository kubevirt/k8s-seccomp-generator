package tracing

import (
	"os"
	"os/exec"
)

type TracingConfiguration struct {
  PodName string `json:"podName"`
}

type Tracer struct { 
  falcoProcess os.Process
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
  // we have to call Process.Release when stopping it
  err := falcoCommand.Start()
  if err != nil {
      return err
  }
  t.falcoProcess = *falcoCommand.Process
  return nil
}

// Stop the tracer by sending interrupt to the Falco process
func (t *Tracer) Stop() error {
  err := t.falcoProcess.Signal(os.Interrupt)
  if err != nil {
    return err
  }
  t.falcoProcess.Release()
  return nil
}
