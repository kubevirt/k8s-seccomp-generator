# Configuring Falco on Kubernetes Cluster

Falco comes with two components, one in the kernel space and the other in the user space. The component in the kernel space is responsible for tracing the syscalls and the component in the user space will use that information to provide the feature set that Falco advocates.

![alt text](https://github.com/kubevirt/k8s-seccomp-generator/blob/main/docs/falco.png?raw=true)

Falco has three different drivers for tracing:
1. eBPF 
2. kernel-module
3. [modern-bpf](https://falco.org/blog/falco-modern-bpf/)

Falco provides a list of pre-compiled kernel modules for some of the common distributions. They can be found [here](https://falcosecurity.github.io/kernel-crawler/?arch=x86_64&target=CentOS). But we cannot depend on the availability of a pre-compiled kernel driver, hence compiling the driver on the node would be the best bet. 

Modern-bpf cannot be used because it requires minimum kernel version `5.8.0`.

For using eBPF, we have to compile and build the Falco eBPF probe on the node.

We can either go with compiling a kernel module or compiling the eBPF probe. We will look at the difference between the two approaches and see which one suits best. 

## Compiling the kernel module on the node

Falco provides a kernel module that we can install on the node. This kernel module will be responsible for tracing the syscalls and sending them to the user space Falco program. 

For installing the kernel module on the node, we would have to follow the following steps (for a centos distribution):

```bash
curl -L -O https://download.falco.org/packages/bin/x86_64/falco-0.35.0-x86_64.tar.gz
tar -xvf falco-0.35.0-x86_64.tar.gz
cp -R falco-0.35.0-x86_64/* /
yum install -y dkms make linux-headers-$(uname -r) 
falco-driver-loader
```

The [falco-driver-loader](https://github.com/falcosecurity/falco/blob/master/scripts/falco-driver-loader) is a script provided by Falco to ease the installation process. It automatically tries to find a pre-built driver for the distribution, and if it not found, it will compile the module and install it using `dkms`. [Dkms](https://linuxhint.com/dkms-linux/), Dynamic Kernel Module Support program/framework allows us to install and load kernel modules that are not a part of the kernel's source tree.

NOTE: For a different distribution, we would have to use the appropriate package manager.

## Compiling the eBPF probe on the node

If we are using an eBPF probe, then we would have to compile the eBPF probe on the node.

The steps would look something like this for centos:

```bash
curl -L -O https://download.falco.org/packages/bin/x86_64/falco-0.35.0-x86_64.tar.gz
tar -xvf falco-0.35.0-x86_64.tar.gz
cp -R falco-0.35.0-x86_64/* /
yum install -y kernel-devel-$(uname -r) clang llvm
falco-driver-loader bpf
```

Here we download Falco and use the `falco-driver-loader` to load a pre-compiled eBPF probe or compile the eBPF probe if not found. The probe can then be used the Falco binary to trace Syscalls.

## Kernel module vs eBPF probe

Although both the methods are viable, we chose eBPF over kernel module because of the following reasons:
1. A faulty kernel module could potentially panic or crash a Linux kernel
2. Loading a kernel module is not always trusted or allowed in some environments
3. The eBPF probe can be dynamically loaded into a kernel at runtime, and does not require using tools like dkms, modprobe, or insmod to load the program

The only argument we can have with using eBPF probe is that we need at least Linux kernel version 4.14, but I think that's a tradeoff we have to make. 

## Automating the installation process

We won't always have the access to ssh into the node and hence we can't ssh into the node and run a script that configures Falco on the node.

We can containerise this approach and run it as a pod. The Dockerfile for such a container image will look like this:

```bash
FROM quay.io/centos/centos:stream8
COPY ./probe_configurer.sh ./installer.sh
RUN curl -L -O https://download.falco.org/packages/bin/x86_64/falco-0.35.0-x86_64.tar.gz
RUN tar -xvf falco-0.35.0-x86_64.tar.gz
RUN chmod +x ./installer.sh
CMD ./installer.sh && cp /root/.falco/falco-bpf.o /usr/_falco/falco-bpf.o
```

where 'probe_configurer.sh' is:

```bash
#!/bin/bash
dnf install -y kernel-devel-$(uname -r)
dnf install -y make clang llvm
mkdir -p /usr/_falco/falco && cp -R falco-0.35.0-x86_64/* /usr/_falco/falco/ && /usr/_falco/falco/usr/bin/falco-driver-loader bpf
```

The pod will look like this:

```yaml
apiVersion: v1
kind: Pod
metadata: 
  name: falco-loader
spec: 
  volumes:
    - name: var
      hostPath:
        path: /var
    - name: usr
      hostPath:
        path: /usr
    - name: lib
      hostPath:
        path: /lib
    - name: etc
      hostPath:
        path: /etc
  containers:
    - name: probe-loader
      image: nithishdev/falco-loader:centos-stream8
      imagePullPolicy: Always
      volumeMounts:
        - name: var
          mountPath: /var
        - name: usr
          mountPath: /usr
        - name: lib
          mountPath: /lib
        - name: etc
          mountPath: /etc
      securityContext: 
        privileged: true
```

We add `/var`, `/usr`, `/etc` and `/lib` as the volume mounts of the Pod. All these mounts are needed so that we can install the `kernel-devel` packages on the node rather than the container. Therefore, we will have different container images for different distributions.

This Pod will configure Falco by compiling the eBPF hook and installing Falco on the node. 
We can use this Pod to install and configure Falco on all the nodes of the cluster. 

