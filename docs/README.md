# Documentation

In order to create an accurate seccomp profile for a pod, we need to know precisely what syscalls the pod will need throughout it's life cycle. We can get this list of syscalls by tracing the pod during it's runtime and mimicking what it will do through functional tests.  

The process of generating seccomp profile for a Kubernetes Pod has the following steps:
1. Auditing the syscalls made by the Pod
2. Using the syscalls list to generate Seccomp profile

There are different tools and approaches for auditing syscalls. 
This section contains documents that review relevant technologies to this project and explain why a particular approach or a tool has been chosen over the others.

This folder contains a set of documents that discuss:
1. [The different approaches available for tracing the syscalls](https://github.com/kubevirt/k8s-seccomp-generator/blob/main/docs/different-approaches-for-auditing-syscalls.md)
2. [The architecture of the k8s-seccomp-generator tool](https://github.com/kubevirt/k8s-seccomp-generator/blob/main/docs/Architecture.md)

It also contains some of the findings from my research on:
1. [Configuring falco on Kubernetes cluster](https://github.com/kubevirt/k8s-seccomp-generator/blob/main/docs/configuring-falco-on-cluster.md)
2. [Tracing syscalls using falco](https://github.com/kubevirt/k8s-seccomp-generator/blob/main/docs/tracing-syscalls-using-falco.md)
