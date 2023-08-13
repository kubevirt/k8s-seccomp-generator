# Kubernetes Seccomp Generator

The Kubernetes Seccomp Generator is a tool designed to simplify the process of generating Seccomp profiles for Kubernetes pods. It provides an intuitive interface for tracing the system calls made by a specific pod and automatically generates a corresponding Seccomp profile based on the observed syscalls.

## Table of Contents

- [Introduction](#introduction)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Introduction

Kuberentes Seccomp Generator uses [Falco](https://falco.org) to trace the syscalls made by the pod and uses that to generate seccomp profile for it. 
Seccomp, short for Secure Computing Mode, is a Linux kernel feature that enables granular control over system calls made by processes. By employing Seccomp profiles, you can significantly reduce the attack surface of your pods by allowing only a specific set of syscalls necessary for their intended functionality.

## Installation

To install the Kubernetes Seccomp Generator, follow these steps:

1. Clone the repository from GitHub:

   ```bash
   git clone https://github.com/sudo-NithishKarthik/kubernetes-seccomp-generator.git
   ```

2. Change into the project directory:
   ```bash
   cd kubernetes-seccomp-generator
   ```
3. Deploy the required components on the cluster:
   ```bash
   scripts/install.sh $nodeos
   ```

   $nodeos here should be the OS distribution your nodes on the cluster uses. This tool currently only supports `centos-stream8`

## Usage

To use the Kubernetes Seccomp Generator, follow these steps:

--TO BE UPDATED--

## Contributing

We welcome contributions to the Kubernetes Seccomp Generator! If you encounter any issues, have suggestions for improvements, or would like to add new features, please feel free to open an issue or submit a pull request on the project's GitHub repository.

## License

