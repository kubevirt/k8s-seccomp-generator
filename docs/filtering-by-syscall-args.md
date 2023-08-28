> This document is a work in progress

# Filtering Syscalls Based on Arguments with Seccomp Profiles

In our effort to enhance security, we aim to filter syscalls based on specific arguments provided. This enables fine-grained control over system calls, allowing us to grant or deny access based on predefined conditions. A practical scenario could involve allowing only certain arguments for the `ioctl` syscall, effectively controlling its behavior on a granular level, such as permitting it for a particular file descriptor (fd).

To achieve this level of control, we can utilize the option to add syscall arguments to the seccomp profile. A seccomp profile could be constructed as follows:

```json
{
    "defaultAction": "SCMP_ACT_ERRNO",
    "defaultErrnoRet": 38,
    "syscalls": [
        {
            "name": "personality",
            "action": "SCMP_ACT_ALLOW",
            "args": [
                {
                    "index": 0,
                    "op": "SCMP_CMP_EQ",
                    "value": 2080505856,
                    "valueTwo": 0
                }
            ]
        }
    ]
}
```

In this example, we've allowed the `personality` syscall with the condition that its first argument should be equal to 2080505856. The seccomp profile offers various operations such as "not equal to," "less than or equal to," and more for filtering based on arguments.

It's worth noting that the usage of arguments in seccomp profiles presents challenges, as highlighted in [this article](https://lwn.net/Articles/822256/). The complexities arise when attempting to filter syscalls based on complex argument types, such as structs.

To explore the feasibility of utilizing syscall arguments for filtering, we've experimented with custom rules that log the syscall argument using Falco. For instance:

```yaml
customRules:
  myrule.yaml: |-
    - rule: myrule
      desc: just for testing purposes
      condition: >
        k8s.pod.name = nginx
      output: >
        (SYSCALL=%syscall.type, args: %evt.args)
      priority: NOTICE
```

From the image below which shows the results from running Falco with the above rule, it can be concluded that Falco can be used to filter syscalls based on arguments.
Along with that, this test also shows Falco's ability to interpret basic syscalls arguments. For instance, the logs shown below have the file descriptor (`/dev/pts/2`) along with the file descriptor number (`fd=2`).

![Test Results](https://github.com/kubevirt/k8s-seccomp-generator/blob/main/docs/filtering-syscalls-args.png?raw=true)


However, for our use case, we aim to filter syscalls based on simple argument types like integers or booleans. Filtering complex arguments like structs could introduce unnecessary complications.

For instance, considering the `ioctl` syscall with arguments:
1. fd (file descriptor) - int
2. request - int
3. argument - untyped pointer

Our immediate goal is to filter syscalls based on the request code argument, keeping the approach straightforward and focused.

One approach to achieve this could be to create a script that processes Falco's output and generates a corresponding seccomp profile. This script would be aware of specific syscall behaviors, such as the need to filter the `ioctl` syscall based on request codes. By integrating this scripted process, we can seamlessly generate seccomp profiles that filter syscalls based on both syscall names and their arguments.
