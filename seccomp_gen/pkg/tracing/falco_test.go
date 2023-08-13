package tracing

import (
	"testing"

	"gopkg.in/yaml.v2"
)

func TestCreateFalcoRule(t *testing.T) {
	tests := []struct {
		name              string
		tracingConf       TracingConfiguration
		expectedCondition string
	}{
		{
			"with Pod name",
			TracingConfiguration{PodName: "virt-launcher"},
			"k8s.pod.name=virt-launcher",
		},
		{
			"with container name",
			TracingConfiguration{ContainerName: "virt-launcher"},
			"container.name=virt-launcher",
		},
		{
			"with Pod label",
			TracingConfiguration{PodLabel: map[string]string{"app.kubevirt.io": "virt-launcher"}},
			"k8s.pod.label.app.kubevirt.io=virt-launcher",
		},
		// {
		//   "with Pod name and container name",
		//   TracingConfiguration{PodName: "virt-launcher", ContainerName: "compute"},
		//   "k8s.pod.name=virt-launcher and k8s.container=compute",
		// },
	}
	for _, tt := range tests {
		res, err := CreateFalcoRule(tt.tracingConf)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		var rule []falcoRule
		err = yaml.Unmarshal(res, &rule)
		if err != nil {
			t.Errorf(err.Error())
			return
		}
		if rule[0].Condition != tt.expectedCondition || rule[0].Output != "Syscall Values: (syscall=%syscall.type)" {
			t.Errorf("got %v", rule)
			return
		}
	}
}
