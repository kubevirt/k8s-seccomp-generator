# Kubernetes Seccomp Generator

The Kubernetes Seccomp Generator is a tool designed to simplify the process of generating [Seccomp](https://man7.org/linux/man-pages/man2/seccomp.2.html) profiles for Kubernetes pods. It provides an intuitive interface for tracing the system calls made by a specific pod and automatically generates a corresponding Seccomp profile based on the observed syscalls.

## Table of Contents

- [Introduction](#introduction)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Introduction

Kuberentes Seccomp Generator uses [Falco](https://falco.org) to trace the syscalls made by the pod and uses that to generate seccomp profile for it. 
Seccomp, short for Secure Computing Mode, is a Linux kernel feature that enables granular control over system calls made by processes. By employing Seccomp profiles, we can significantly reduce the attack surface of your pods by allowing only a specific set of syscalls necessary for their intended functionality.

## What does it solve?

[Seccomp](https://man7.org/linux/man-pages/man2/seccomp.2.html), short for "secure computing mode," is a security feature in Linux operating systems that provides a way to restrict the system calls that a process can make. By limiting the set of allowed system calls, seccomp can help reduce the potential attack surface and improve the security of a system.

Container engines like [Docker](https://www.docker.com/), [Podman](https://podman.io/) etc. allow us to limit the set of allowed syscalls by [applying a seccomp profile to a container](https://docs.docker.com/engine/security/seccomp/).  Container engines offer their own default seccomp profile. However, we cannot assume that one size fits all. Therefore, the default profile may either permit syscalls that are in fact not required by the workload or prohibit legitimate syscalls. 

Kubernetes too allows us to apply Seccomp profiles loaded onto the node to the Pods and Containers - https://kubernetes.io/docs/tutorials/security/seccomp/ 

In order to create an accurate seccomp profile for a pod, we need to know precisely what syscalls the pod will need throughout it's life cycle. We can get this list of syscalls by tracing the pod during it's runtime and mimicking what it will do through functional tests. This tool simplifies this process by providing an easy-to-use interface to generate seccomp profile for a Kubernetes pod.  

## Installation

Install the secgen cli tool.

--TO BE UPDATED--
 
## Usage

To use the Kubernetes Seccomp Generator, follow these steps:

### 1. Install required components on the cluster

`secgen install $node-os`

This configures the node and installs the required components on the cluster.
$node-os is the distribution of the operating system used on the node of the kuberentes cluster. Currently we only support `centos-stream8`, but this can be extended easily.

### 2. Start tracing

`secgen trace start $selector`

This starts tracing the pod referred by the $selector.
$selector can be 'pod.name=$name', 'container.name=$name' or 'pod.label.$label=$value'

### 3. Stop tracing

`secgen trace stop`

This stops tracing and outputs the seccomp profile generated for the pod. 

NOTE: Only single node clusters are supported. Support for multiple nodes will be added later.

## Contributing

We welcome contributions to the Kubernetes Seccomp Generator! If you encounter any issues, have suggestions for improvements, or would like to add new features, please feel free to open an issue or submit a pull request on the project's GitHub repository.

Current progress and TODOs are being tracked here: https://github.com/kubevirt/k8s-seccomp-generator/issues/2

## License

Apache License
