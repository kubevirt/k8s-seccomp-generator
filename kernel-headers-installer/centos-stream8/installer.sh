dnf install -y kernel-devel-$(uname -r)
dnf install -y make clang llvm
falco-driver-loader bpf
