package tracing

import (
	"testing"
	"gopkg.in/yaml.v2"
)


func TestCreateFalcoRule(t *testing.T){
    tracingConf := TracingConfiguration{PodName: "virt-launcher"}
    res,_ := CreateFalcoRule(tracingConf)
    var rule falcoRule
    yaml.Unmarshal(res, &rule)
    if rule.Condition != "k8s.pod.name=virt-launcher" || rule.Output != "Syscall Values: (syscall=%syscall.type)" {
        t.Errorf("got %v", rule)
  }
}
