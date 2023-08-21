# Comprehensive System Architecture

The underlying architecture of this system design encompasses two pivotal components, each playing a distinct role in achieving the desired functionality. These components are as follows:

1. **In-Cluster Component (Component A)**  
   This element operates within the cluster environment, ensuring the seamless execution of various tasks associated with the system's operation.

   A visual representation of the deployment plan is illustrated below:
   ![Deployment Plan](https://github.com/kubevirt/k8s-seccomp-generator/blob/main/docs/DeploymentPlan.png?raw=true)

   Component A, meticulously engineered to run one instance per node, is responsible for orchestrating the core functions of the system. Its responsibilities are subdivided into three interrelated tasks:

   1. **Setting up Falco**: Component A assumes the role of configuring and overseeing the Falco intrusion detection system's deployment. It meticulously manages all aspects related to Falco's operation within the cluster.
   2. **Parsing the Logs**: Component A leverages a dedicated binary in conjunction with [Falco Alert Channels](https://falco.org/docs/alerts/channels/) to effectively parse the system logs. These logs are then meticulously stored within a structured `data.json` file.
   3. **Generating the Seccomp Profile**: Component A further extends its capabilities to encompass the generation of Seccomp profiles. These profiles are vital for enhancing security by defining permissible syscall behaviors.

   To enable these functionalities, Component A integrates a Falco instance, an API server, and the ability to process and store logs, ultimately facilitating the seamless operation of the entire system.

2. **CLI Client**  
   The Command-Line Interface (CLI) client serves as the orchestrator of the system. It acts as the central control point for various operational aspects, wielding a pivotal role in executing and managing system-wide functionalities.

   Distinguished as the "smart" component, the CLI client acts as the decision-maker, directing Component A's actions. Component A, in contrast, assumes a "dumb" role, following the instructions imparted by the CLI client. The CLI client's responsibilities encompass tasks such as:

   - **Merging Syscalls**: The CLI client adeptly merges syscalls originating from diverse nodes, thereby fostering a cohesive view of system behavior and potential vulnerabilities.
   - **Seccomp Profile Generation**: It governs the generation of Seccomp profiles, a critical security aspect that restricts syscall activities to enhance system integrity.
   - **System-Level Configuration Management**: The CLI client actively manages system-level configurations, including user-defined preferences such as pods to trace and the inclusion of syscall arguments.

## User-Centric Workflow

As an end user, engaging with this system involves a structured workflow that unfolds as follows:

1. **Installation and Deployment**  
   The initial step entails the installation of the CLI client, which subsequently serves as the primary interface for interacting with the system. The CLI client plays a pivotal role in deploying the entire system within a Kubernetes cluster.

   Installation may include the prerequisite of deploying kernel headers onto hosts, ensuring a foundational environment conducive to the system's operation. The CLI client effectively utilizes kubectl manifests or helm charts to orchestrate the deployment process.

2. **Component A Initialization**  
   Upon successful deployment, Component A springs to life within designated cluster nodes. Its primary tasks involve ensuring the seamless integration and readiness of the Falco eBPF driver, a crucial element for the subsequent execution of Falco. Should any issues arise during this phase, Component A diligently updates the associated Kubernetes Object to reflect relevant errors.

3. **Configuration Customization**  
   As a user, you are empowered to tailor the system's behavior to align with your requirements. This entails configuring options such as specifying the pods to be traced and determining the inclusion of syscall arguments in generated profiles.

4. **Commencement of Tracing**  
   With configurations set, you command the CLI client to initiate the tracing process. In response, Component B, a module within the CLI client, engages in a thorough health assessment of all Component A instances. Subsequently, it dispatches requests to the respective API servers hosted by Component A instances, effectively kickstarting the tracing operation. Each Component A instance uses the supplied parameters to commence Falco's operation, parsing and storing logs within designated `data.json` files.

5. **Tracing Termination and Profile Generation**  
   When you opt to conclude the tracing process, Component B orchestrates the cessation of tracing activities across all Component A instances. Subsequently, it diligently collects the generated `data.json` files from each instance. Armed with this data, the CLI client merges the information, harmonizing insights derived from diverse nodes. With this cohesive dataset, the CLI client generates Seccomp profiles that encapsulate syscall behaviors. The profiles reflect the user-defined configuration, contributing to bolstered system security.

In essence, the system architecture adeptly amalgamates Component A's intricate orchestration capabilities with the CLI client's intelligent management, culminating in a comprehensive solution that empowers users to enhance security and operational integrity within their Kubernetes clusters.
