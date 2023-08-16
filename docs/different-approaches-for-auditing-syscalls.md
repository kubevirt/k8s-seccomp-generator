# Different Approaches for Tracing Syscalls

Various methods can be employed to trace the syscalls initiated by a Kubernetes Pod. These approaches encompass:

1. Utilizing the Kubernetes Seccomp profile feature (with the Log action)
2. Employing tools like Strace, Kstrace, or similar alternatives
3. Harnessing the power of eBPF (Extended Berkeley Packet Filter)
4. Leveraging tools such as Falco or Sysdig

## Leveraging the Kubernetes Seccomp Profile Feature

A seccomp profile can be crafted, employing the LOG action as its default behavior. The LOG action ensures the recording of all syscalls executed by the pod. This mechanism essentially relies on the [`seccomp-bpf`](https://www.kernel.org/doc/html/v4.16/userspace-api/seccomp_filter.html) filter mode, utilizing the [`prctl`](https://man7.org/linux/man-pages/man2/prctl.2.html) system call under the hood. This profile can be applied to the designated Kubernetes pod.

The resulting syscall logs will be stored either within the `auditd` log file (/var/log/audit/audit.log) or the `syslog` file, contingent upon system configuration. If the system is running the Linux audit daemon, logs will default to the audit.log file; otherwise, they will be directed to the syslog file.

This approach is relatively straightforward and obviates the necessity of writing and maintaining additional code.

However, it's important to note the drawbacks associated with this approach:

1. Inability to distinguish between syscalls from different pods:
   - This limitation becomes apparent when other pods possess applied seccomp profiles, making it challenging to discern syscalls from distinct pods. While this isn't an insurmountable issue when tracing a single pod, multiple concurrently traced pods with seccomp profiles would require profile modifications to cease logging and facilitate proper differentiation.
2. Limited information in logs:
   - The log entries lack syscall arguments, impeding precision. A more granular seccomp profile could be achieved by incorporating syscall argument details, enabling fine-tuned control over syscall execution based on parameter usage.
3. Log duplication:
   - In our use case, the frequency of syscall occurrences is of lesser concern than their existence. Applying a seccomp profile with LOG as the default action yields duplicate entries for syscalls, leading to log inflation.

## Utilizing Strace, Kstrace, or Similar Alternatives

[strace](https://man7.org/linux/man-pages/man1/strace.1.html) stands as a versatile command-line utility prevalent in Unix-like systems, particularly Linux. It facilitates the tracing and analysis of system calls and signals invoked during program execution. By delving into the interactions between processes and the operating system kernel, strace provides insights into system behavior.

The following methods are variations predicated on strace.

### Kstrace

[Kstrace](https://github.com/MichaelWasher/kstrace) entails tracing pod syscalls through a privileged pod linked to the target pod. Strace is utilized within the target pod to track system calls. However, a notable drawback arises from the synchronization intricacies between pod launch and strace's initiation. These synchronization challenges can introduce race conditions and potential syscall loss.

Alternatively, in scenarios where kstrace isn't employed, access to the Kubernetes node enables the use of `crictl`, a command-line utility for CRI-compatible container runtimes. With this approach, the container's PID can be retrieved, permitting strace to trace the syscalls of the respective container. Nonetheless, similar to kstrace, initiating tracing immediately upon container start remains uncertain. This limitation extends to analogous tools.

### Wrapping the Container Entrypoint with `strace`

An alternative avenue involves encapsulating the entry point command in a `strace` wrapper within the Docker image. However, feasibility concerns may surface, and seamless integration with test suites could prove intricate.

### Employing OCI Hooks

[OCI hooks](https://man.archlinux.org/man/oci-hooks.5.en) offer a method to configure hooks for Open Container Initiative containers, activating exclusively for containers necessitating their functionalities and pertinent stages. This mechanism facilitates integration with the container's lifecycle for executing specific processes.

Combining tracing tools with OCI hooks is a promising option. Configuring a binary during the `preStart` stage of the container via OCI hooks can initiate strace (or similar tools) for the container process. The OCI hook's input data, including PID, empowers the binary's operation. Strace output files are subsequently stored and retrieved. To circumvent blocking, a child process initiates strace, ensuring the container's unimpeded commencement.

With refinement, this approach can be optimized further. Notably, the [–seccomp-bpf flag](https://pchaigno.github.io/strace/2019/10/02/introducing-strace-seccomp-bpf.html) could yield transparent performance enhancements.

Effective utilization of OCI hooks necessitates scripting to configure hooks within the Kubernetes cluster.

## Harnessing eBPF

[eBPF](https://ebpf.io/) constitutes a sophisticated Linux technology, enabling secure code injection into the kernel for tasks such as tracing, filtering, and monitoring. It empowers dynamic tracing, network packet filtering, and performance analysis, all without altering the kernel source. eBPF programs, written in a subset of C, are loaded into the kernel, providing real-time insights into system behavior and enhancing security.

From the perspective of the Linux kernel, `containers` are synonymous with processes. By obtaining a container process's PID namespace upon initiation, an eBPF program can be compiled and loaded to interface with the `sys_enter` tracepoint. This approach facilitates logging of syscalls executed solely within the container's PID namespace.

However, this eBPF program must be loaded just prior to the container's launch. OCI hooks offer a viable solution. The `preStart` hook, executed before container commencement, can activate a binary facilitating eBPF program loading.

The creation of an OCI hook within CRI-O involves configuring hook behavior, specifying the binary and triggering conditions. The existing [oci-seccomp-bpf-hook](https://github.com/containers/oci-seccomp-bpf-hook) implements a similar approach for generating seccomp profiles.

It's important to note certain prerequisites for this hook's functionality:

1. Root privileges (CAP_SYS_ADMIN) are essential; rootless containers are incompatible.
2. The binary must possess CAP_SYS_ADMIN privileges. This won't impede rootless VMs.
3. The presence of the bcc toolchain and kernel headers on the node is necessary for compiling and loading BPF programs.

The oci-seccomp-bpf hook serves container runtimes, like runc, in applying seccomp filters. However, this hook setup varies, contingent on the container runtime and engine. Consequently, code implementation for our Kubernetes cluster's hook configuration is indispensable.

A limitation lies in OCI hooks' compatibility solely with CRI-O, lacking support for containerd.

## Employing Falco

[Falco](https://falco.org) emerges as a cloud-native security tool tailored for Linux systems. It leverages custom rules applied to kernel events, enriched with container and Kubernetes metadata, to furnish real-time alerts.

Falco supports tracing via a kernel module, eBPF, and modern bpf. Kernel space components undertake tracing, with user space counterparts harnessing the information for practical features.

Efficiently tracing syscalls within a Kubernetes Pod remains a priority. Yet, Falco's deployment across clusters poses a generalized challenge. Although a [helm chart](https://github.com/falcosecurity/charts) facilitates Kubernetes cluster installation, issues persist, as outlined in [configuring-falco-on-cluster](https://github.com/kubevirt/k8s-seccomp-generator/blob/main/docs/configuring-falco-on-cluster.md).

Comparing the Falco and bpf-hook approaches reveals:

1. The bpf-hook approach's challenge lies in widespread deployment across Kubernetes clusters.
2. OCI hooks' limited support, e.g., exclusion from containerd, presents an obstacle.
3. Falco boasts a well-tested solution, offering robust syscall tracing for Kubernetes Pods. It features syscall argument support and excels in flexibility, irrespective of the container engine.

Hence, Falco emerges as the optimal solution for our specified use case.

