dnf install -y kernel-devel-$(uname -r)
dnf install -y make clang llvm
mkdir -p /usr/_falco/ && cp -R falco-0.35.0-x86_64/* /usr/_falco/falco/ && /usr/_falco/falco/usr/bin/falco-driver-loader bpf
