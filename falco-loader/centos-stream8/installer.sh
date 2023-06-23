dnf install -y kernel-devel-$(uname -r)
dnf install -y make clang llvm
/usr/_falco/falco/usr/bin/falco-driver-loader bpf
