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
  condition := "k8s.pod.name="+t.PodName
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
