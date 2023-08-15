# Different Approaches for Auditing Syscalls

There are different approaches we can use to trace the syscalls:
1. Kubernetes Seccomp profile (using the Log action)
2. Using Strace, Kstrace or similar alternatives
3. Using eBPF
4. Using Falco or Sysdig


## Using Kubernetes Seccomp Profile Feature

We can create a seccomp profile that uses the LOG action as the default action. This under the hood uses the `seccomp-bpf` filter mode with `prctl()`.  This profile can be applied to the Kubernetes pod.

The logs, which contain the list of syscalls made by the pod, will be stored either in the `auditd` log file (/var/log/audit/audit.log) or the `syslog` file depending on the system configuration. If the system has linux audit daemon running, the logs will be stored in the audit.log file by default, else it will be stored in the syslog file.

In this approach, we won’t have to write or maintain any code.

Some of the disadvantages of using this approach are:
1. We cannot differentiate between the syscalls made by different pods 
	- This will become an issue when other pods have seccomp profiles applied as we won’t be able to differentiate between the syscalls made by different pods. This will not be a blocker as long as we can ensure that there is only one pod being traced at the moment. If there are other pods that have a seccomp profile applied, we would have to do some extra work of modifying the profiles to stop logging in order for this approach to work. 
2. Less information in the logs
	- For instance, the logs don’t list the arguments of the syscall which is a limitation for us. 
	- We cannot control what information is exposed in the logs 
3. Duplication in the logs 
	- For our use case, we are not really concerned about how many times a syscall is made by a pod. We are only concerned about whether or not the pod makes a particular syscall. When we apply a seccomp profile with LOG as the default action, all the syscalls made by the pod will be added in the log file which will contain a lot of duplicates. This increases the size of the log file. I had observed that the logs from a single VirtLauncher pod is around 24 MBs. This issue magnifies when we have multiple VirtLauncher pods running.
 4. We also don't have much control over how these things are being logged. For us, the whole logging mechanism is a black box and it will be difficult for us to debug if something does not work.

## Using Strace, Kstrace or Similar Alternatives

### Kstrace

If we take the [kstrace](https://github.com/MichaelWasher/kstrace) tool (which is a wrapper around strace), it traces the syscalls made by a pod by creating a privileged pod and attaching it to the target pod, which will use strace to trace the system calls of the target pod. But the problem with kstrace is that, we cannot guarantee that kstrace will start tracing simultaneously when the target pod starts, and hence we might lose some syscalls.

  
If we do not want to use kstrace, if we have access to the kubernetes node, we can use `crictl` to get the PID of the container (since every container is just a linux process at the end of the day), with which we can run strace to trace the syscalls the container is making. But the disadvantage with this approach is that we still cannot guarantee that we can start tracing as soon as the container starts. The same argument goes for other similar tools as well. 

### Wrapping the container entrypoint with `strace`

Another approach would be to modify the docker image to wrap the entry point command in `strace`. But it will not always be feasible to modify the docker image and it will be complex to integrate this with tests suites.

### Using  OCI hooks

OCI hooks provide a way for users to configure the intended hooks for Open Container Initiative containers so they will only be executed for containers that need their functionality, and then only for the stages where they're needed. Learn more about OCI hooks [here](https://man.archlinux.org/man/oci-hooks.5.en).

Using these tracing tools along with OCI hooks is a viable option. We can have a binary configured to the  `preStart` hook of the container lifecycle. This binary will start `strace` (or any other similar tools) for the container process (using the PID which we get from the stdin). Strace output files can be stored in some location and can be retrieved later. A new child process will be created by the `preStart` binary which will start `strace`. This is because if the `preStart` binary is blocking, then the container will not start. 

But the problem with this approach is finding when to stop tracing. A trivial solution would be to use `(while kill -0 $pid; do sleep 1; done)` where $pid will be the pid of the container process. But this is inefficient.

There is no direct way to get notified about the exit/completion of a process. Doing a periodic scan of the /proc/pseudo file system is another approach but it has the same problems we have with the signal 0 kill approach and doesn’t make it any better. We cannot use the `postStop` OCI hook since we won’t have access to the `preStart` hook binary’s child process to signal it to stop tracing. 

With some amount of engineering effort, we can make the other aspects of this approach better. For instance, we can use the [–seccomp-bpf flag](https://pchaigno.github.io/strace/2019/10/02/introducing-strace-seccomp-bpf.html) to make it faster. And we can make strace stream the logs to a child process (by child process, I mean that this process will be a child of the OCI `preStart` hook binary, this has to be a child process because the container will not start unless the binary is non-blocking) before writing it to avoid duplicate logs.

When we are using OCI hooks, we would have to write and maintain code (it might just be a shell script) that sets up the hook in our kubernetes cluster.

## Using eBPF

As long as the linux kernel is concerned, `containers` are just processes. When the container starts, if we can get the PID namespace of the container process, then we can compile and load an eBPF program that hooks into the `sys_enter` tracepoint and logs only the syscalls made by the processes that are within the PID namespace of the container.  

But, the eBPF program has to be loaded just before the container starts. We can use the OCI hooks for this. OCI hooks allow us to run binaries at different life cycles of a container.  We can use the `preStart` hook which will run a binary when the container namespaces are created, but the container has not started yet.

We can create an OCI hook in CRI-O by creating a hook configuration that describes what binary to be called and when it should be triggered ([https://cloud.redhat.com/blog/extending-the-runtime-functionality](https://cloud.redhat.com/blog/extending-the-runtime-functionality) for more information). There is already an OCI hook ([oci-seccomp-bpf-hook](https://github.com/containers/oci-seccomp-bpf-hook)) that generates a seccomp profile for a container by using the same approach.

NOTE: Kubernetes also provides a way to hook into the [container lifecycles](https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/) but it is not good enough for us since the `postStart` hook does not guarantee that it will be run before the entrypoint of the container.

We would have to set up the hook to start the oci-seccomp-bpf binary in our node. This would require the node to have bcc toolchain and kernel headers. The hook will be configured in such a way that it will run for all the containers with the `io.containers.trace-syscall` annotation. This would mean that we would have to set the hook annotation for the VirtLauncher pod while running tests. Once the container exits, the profile will be outputted to a specific location. Note that if there are multiple containers running in a pod, we have to merge the whitelisted syscalls for all the containers to arrive at the final seccomp profile for the pod.

Some of the requirements for this hook to work is: 
1. Root privileges (CAP_SYS_ADMIN) are needed.  Hook will not work with rootless containers.
2. The binary needs to have CAP_SYS_ADMIN to run. We can still start rootless VMs and will not be a problem for us. 
3. The bcc tool chain and kernel-headers should be present in the node to be able to compile and load BPF programs
    

The oci-seccomp-bpf hook is used by container runtimes, such as runc, to apply seccomp filters to container processes. There is no generalised solution that automates the process of setting up this hook in kubernetes, most probably because it is dependent on the container runtime and container engine used. Therefore, we would have to write and maintain code to set up this hook for our kubernetes cluster.

The main problem with this approach is the support for different container runtimes. For instance, containerd does not have support for OCI hooks, and hence we can't use it.

## Using Falco

[Falco](https://falco.org) is a cloud-native security tool designed for Linux systems. It employs custom rules on kernel events, which are enriched with container and Kubernetes metadata, to provide real-time alerts.

![alt text](https://github.com/kubevirt/k8s-seccomp-generator/blob/main/docs/falco.png?raw=true)

Falco supports tracing using kernel module, eBPF and modern bpf. Tracing is done by the component in the kernel space and the user space components make use of that information to provide useful features. 

For us, we need to be able to trace syscalls of a particular Kubernetes Pod. 

Comparing the falco and bpf-hook approach, we can say that:
1. In case of the bpf-hook approach, the hurdle is with deployment. Also this approach is not supported by some container engines. 
2. Whereas in case of the Falco approach, the Falco project already has a well-tested solution for tracing Syscalls for Kubernetes Pods. It supports syscall arguments and is very flexible. It is also independent on the container engine being used. If we can find a generalised solution for installing Falco on the cluster, then this would be the ideal solution for us. 
