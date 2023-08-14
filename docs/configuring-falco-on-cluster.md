In order for the Falco approach to work, we need to make sure that we can install Falco on most distributions without problems. 

Falco provides three different tracing solutions:
1. eBPF 
2. kernel-module
3. modern-bpf

Falco provides a list of pre-compiled kernel modules for some of the common distributions. They can be found [here](https://falcosecurity.github.io/kernel-crawler/?arch=x86_64&target=CentOS). 
But we cannot use it because there are a lot of distributions which they don't have drivers for. 

Modern-bpf cannot be used because it requires minimum kernel version `5.8.0`.

For using eBPF, we have to compile and build the Falco eBPF probe on the fly. 

The steps would look something like this for centos:

```bash
curl -L -O https://download.falco.org/packages/bin/x86_64/falco-0.35.0-x86_64.tar.gz
tar -xvf falco-0.35.0-x86_64.tar.gz
cp -R falco-0.35.0-x86_64/* /
yum install kernel-devel-$(uname -r)
# this needed to compile the eBPF probe on the fly
yum install -y make clang llvm
falco-driver-loader bpf
```

Here we download Falco and use the `falco-driver-loader` to compile the eBPF probe. The probe can then be used the Falco binary to trace Syscalls.

Now we need to automate this approach. First thing that comes up is to ssh into the node and do these things. We can write a script that automates this. But the problem with this approach is that we cannot guarantee that we will have access to the node at all times. 

We can containerise this approach and run it as a pod. 

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

We add `/var`, `/usr`, `/etc` and `/lib` as the volume mounts of the Pod. This Pod will configure Falco by compiling the eBPF hook and installing Falco on the node. 
We can use this Pod to install and configure Falco on all the nodes of the cluster. 

We have to have different container images for different distributions. 

