# Configuration of Falco within the Kubernetes Cluster

When it comes to setting up Falco in a Kubernetes cluster, there are two key components involved: one operating in the kernel space and the other in the user space. The kernel component is responsible for tracing syscalls, while the user space component leverages this traced data to offer the feature set endorsed by Falco.

![alt text](https://github.com/kubevirt/k8s-seccomp-generator/blob/main/docs/falco.png?raw=true)

Falco presents three distinct tracing drivers:
1. [eBPF](https://ebpf.io/)
2. kernel-module
3. [modern-bpf](https://falco.org/blog/falco-modern-bpf/)

It's important to note that adopting modern-bpf necessitates a minimum kernel version of `5.8.0`. This requirement might potentially limit compatibility with the tool (k8s-seccomp-generator); hence, this option is not the most suitable at this moment. However, the possibility of adding modern-bpf support remains open for future consideration.

For utilizing the eBPF approach, the Falco eBPF probe must be compiled and constructed on the node.

In this regard, two pathways are available: compiling a kernel module on the node or compiling the eBPF probe. A comprehensive comparison between these approaches will help us determine the most optimal choice moving forward.

## Compiling the Kernel Module on the Node

Falco provides a traditional kernel module, which serves as a classic method of obtaining the requisite data stream from the kernel. This module's function is to trace syscalls and subsequently transmit them to the Falco user space program.

For CentOS distributions, the installation procedure unfolds as follows:

```
curl -L -O https://download.falco.org/packages/bin/x86_64/falco-0.35.0-x86_64.tar.gz
tar -xvf falco-0.35.0-x86_64.tar.gz
cp -R falco-0.35.0-x86_64/* /
yum install -y dkms make linux-headers-$(uname -r) 
falco-driver-loader
```

To simplify the installation process, the [falco-driver-loader](https://github.com/falcosecurity/falco/blob/master/scripts/falco-driver-loader) script is provided. This script streamlines the installation by searching for a pre-built driver specific to the distribution. In instances where such a driver is not found, the script takes on the task of compiling the necessary module using `dkms`.

## Compiling the eBPF Probe on the Node

When opting for the eBPF probe approach, compiling the probe on the node becomes a necessary step.

For CentOS distributions, this process involves the following steps:

```
curl -L -O https://download.falco.org/packages/bin/x86_64/falco-0.35.0-x86_64.tar.gz
tar -xvf falco-0.35.0-x86_64.tar.gz
cp -R falco-0.35.0-x86_64/* /
yum install -y kernel-devel-$(uname -r) clang llvm
falco-driver-loader bpf
```

In this context, the `falco-driver-loader` script is employed to either load an existing pre-compiled eBPF probe or, if unavailable, compile the required eBPF probe. Once in place, the Falco binary can then leverage this probe to trace syscalls.

## eBPF Probe vs. Kernel Module

Although both methods hold promise, the preference tilts towards eBPF for several reasons:
1. Kernel module malfunctions could lead to Linux kernel crashes or panics.
2. In certain environments, the loading of kernel modules is not universally trusted or permitted.
3. eBPF probes can be dynamically loaded into the kernel during runtime, circumventing the need for tools like dkms, modprobe, or insmod for loading the program.

The tradeoff comes in the form of eBPF's requirement of at least Linux kernel version 4.14. This stipulation, however, is deemed reasonable and acceptable.

## Streamlining Installation

While Falco offers a Helm chart for installation, this option might not be the most suitable for our specific use case. Challenges arise when pre-built eBPF probes or kernel modules are missing. In such instances, manual installation becomes essential, or the Helm chart requires packages to be installed for successful compilation and loading. For instance, the eBPF approach necessitates the installation of the `kernel-modules` package on the node for the successful installation of Falco via the Helm chart. Consequently, a more pragmatic alternative is to adopt our own containerized solution.

A Docker container, known as 'falco-loader,' facilitates the execution of required commands. This container leverages the [`falco-driver-loader`](https://github.com/falcosecurity/falco/blob/master/scripts/falco-driver-loader) script to compile the eBPF hook. The project maintains a collection of OS-specific container images, accessible [here](https://github.com/kubevirt/k8s-seccomp-generator/blob/main/install/falco_loader). This assortment caters to different OS distributions, each of which employs distinct package managers. Presented below is a Dockerfile example targeting CentOS Stream 8:

```Dockerfile
FROM quay.io/centos/centos:stream8
COPY ./probe_configurer.sh ./installer.sh
RUN curl -L -O https://download.falco.org/packages/bin/x86_64/falco-0.35.0-x86_64.tar.gz
RUN tar -xvf falco-0.35.0-x86_64.tar.gz
RUN chmod +x ./installer.sh
CMD ./installer.sh && cp /root/.falco/falco-bpf.o /usr/_falco/falco-bpf.o
```

The accompanying 'probe_configurer.sh' script is as follows:

```
#!/bin/bash
dnf install -y kernel-devel-$(uname -r)
dnf install -y make clang llvm
mkdir -p /usr/_falco/falco && cp -R falco-0.35.0-x86_64/* /usr/_falco/falco/ && /usr/_falco/falco/usr/bin/falco-driver-loader bpf
```

To run this solution as a Kubernetes pod, key directories from the host file system (such as `/var`, `/usr`, `/lib`, and `/etc`) are mounted within the pod. This enables the direct installation of `kernel-modules` on the node, ensuring compatibility. The pod encompasses a privileged container ('probe-loader') responsible for configuring Falco by compiling the eBPF hook and conducting Falco installation on the node.

This pod offers a feasible approach for configuring Falco across all nodes within the cluster.
