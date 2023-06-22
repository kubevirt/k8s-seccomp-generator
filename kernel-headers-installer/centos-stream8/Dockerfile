FROM quay.io/centos/centos:stream8
RUN dnf install -y kernel-devel-$(uname -r)
RUN dnf install -y make clang llvm
RUN curl -L -O https://download.falco.org/packages/bin/x86_64/falco-0.35.0-x86_64.tar.gz
RUN tar -xvf falco-0.35.0-x86_64.tar.gz
RUN cp -R falco-0.35.0-x86_64/* /
RUN falco-driver-loader bpf
CMD mkdir /usr/_falco-probe/ && cp /root/.falco/falco-bpf.o /usr/_falco-probe/falco-bpf.o


