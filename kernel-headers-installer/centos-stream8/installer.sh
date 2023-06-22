dnf install -y kernel-devel-$(uname -r)
dnf install -y make clang llvm
falco-0.35.0-x86_64/usr/bin/falco-driver-loader bpf
