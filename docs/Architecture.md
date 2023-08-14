
This system design consists of two components: 
1. Component A (In-Cluster Component)
2. CLI Client


![alt text](https://github.com/kubevirt/k8s-seccomp-generator/blob/main/docs/DeploymentPlan.png?raw=true)

Component A will be running as _1 per node_. 

We can split the whole procedure into three tasks:
1. Setting up Falco
2. Parsing the logs
3. Generating the Seccomp profile

### Component A

Component A will make sure that Falco is properly configured and running. It is responsible managing Falco. As far as outside world (anything outside of Component A) is concerned, there is no such thing as Falco, it is the Component A that traces and provides the syscalls. 

Component A will have Falco running. It will make use of a binary along with the [Falco Alert Channels](https://falco.org/docs/alerts/channels/) to parse the logs and store them in a `data.json` file. 

Component A will also have an API server running. This will be point of communication for Component A and the CLI client.

### CLI Client

The CLI Client will be like an orchestrator.

This is the _smart_ component of the system, meaning that Component A will be _dumb_ and just does what the CLI says. The CLI client will take care of things like merging syscalls from different nodes, generating a Seccomp profile, storing system level configurations (user configurations like pods to trace, whether syscall arguments should be included, etc. )

## User Flow

As a user, I would first install the CLI and use the CLI to deploy the system in my Kubernetes cluster. 

The CLI will first install the kernel headers on the hosts (or if it is complex, then we can make this a prerequisite). Then the ClI will use the kubectl manifests or helm charts to install our deployments.

Now, as soon as Component A is deployed on a node, it will make sure that the Falco eBPF driver is loaded and Falco is ready to be run. If there are any issues with this, it will update the Kubernetes Object with the errors.

Then, as a user, I will set the configuration options like which pods to trace, whether I need profiles with syscall arguments, etc.

Now, I will ask the CLI client to start tracing. At this point, Component B will first check the health of all Component 'A's and send a request to their API server to start tracing with necessary configuration parameters. 
Component A will then use the parameters to start Falco with the binary to parse and store the logs in the `data.json` file. 

Now, I will ask the CLI to stop tracing and get the Seccomp profile. Now Component B will ask all the 'A's to stop tracing and then get the `data.json` from all the 'A's. After this, the client will merge the data and generate the Seccomp profile, with the given configuration.
