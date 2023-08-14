We can use the Falco helm chart to install Falco on the cluster. The cluster being used here has a single node and the node has centos-stream8.

Source: https://falco.org/docs/getting-started/installation/#centos-rhel

1. In the node
```bash
rpm --import https://falco.org/repo/falcosecurity-packages.asc
yum install -y epel-release
yum install -y dkms 
yum install -y kernel-devel-$(uname -r)
```

2. In the host where kubernetes is running
```bash
helm repo add falcosecurity https://falcosecurity.github.io/charts
helm repo update
helm install falco falcosecurity/falco --namespace falco --create-namespace
```


Result: 
![alt text](https://github.com/kubevirt/k8s-seccomp-generator/blob/main/docs/result-falco-syscalls.png?raw=true)

Now to test out that it works, We can start an `nginx` container and exec into it to see if falco detects anything.

`kubectl logs -f <falco-pod>` 

Log buffering is enabled by default in the falco helm chart (https://github.com/falcosecurity/charts/tree/master/falco#enabling-real-time-logs).
`helm install falco falcosecurity/falco --set tty=true` disabled log buffering and allows us to see logs in real time. 

We create a custom falco rule that allows us to filter for syscalls made by the nginx pod: 

```yaml
customRules:
  myrule.yaml: |-
    - rule: myrule
      desc: just for testing purposes
      condition: >
        k8s.pod.name = nginx
      output: >
        (SYSCALL=%syscall.type)
      priority: NOTICE
```


We do `helm install falco -f custom-rules.yaml falcosecurity/falco --set tty=true`.

Then once the falco pod has started, We start the nginx pod. We can see all the syscalls made by the pod in the falco pod logs.

Output:
![alt text](https://github.com/kubevirt/k8s-seccomp-generator/blob/main/docs/output-falco-syscalls.png?raw=true)

But the problem is, as we can see in the screenshot, there are some places where we have `SYSCALL=<NA>`. This is because falco by default limits the supported list of syscalls because of performance reasons. So, in order for us to get these events as well, we have to run falco with -A flag.

After changing the `%syscall.type` to `%evt.type` along with the `-A` argument, we can get all the syscalls.

Finalised falco rule:

```yaml
customRules:
  myrule.yaml: |-
    - rule: myrule
      desc: just for testing purposes
      condition: >
        k8s.pod.name = nginx
      output: >
        (SYSCALL=%evt.type)
      priority: NOTICE
```

This is because when we use `evt.type`, we will get the logs for all the events, whereas when we use `syscall.type` we will only get the value if it is a syscall. This is the reason why some of them are shown as NA. 

