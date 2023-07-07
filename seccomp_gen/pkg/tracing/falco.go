package tracing

import (
	"gopkg.in/yaml.v2"
)

type falcoRule struct {
  Rule string `yaml:"rule"`
  Desc string `yaml:"desc"`
  Condition string `yaml:"condition"`
  Output string `yaml:"output"`
  Priority string `yaml:"priority"`
}

// CreateFalcoRule generates a Falco rule from the TracingConfiguration
func CreateFalcoRule(t TracingConfiguration) ([]byte, error) {
  var condition string
  if t.PodName != "" {
  condition = "k8s.pod.name="+t.PodName
  }
  if t.ContainerName != ""{
  condition = "container.name="+t.ContainerName
  }
  if t.PodLabel != nil {
    var lbl, val string
    for k,v := range t.PodLabel {
      lbl = k
      val = v
    }
    condition ="k8s.pod.label."+lbl+"="+val
  }
  rules := make([]falcoRule, 0)
  rules = append(rules, falcoRule{
    Rule: "ksecgenRule",
    Desc: "Testing Rule",
    Condition: condition,
    Output: "Syscall Values: (syscall=%syscall.type)",
    Priority: "WARNING",
  })
  res, err := yaml.Marshal(rules)
  if err != nil {
    return nil, err
  }
  return res, nil
}
