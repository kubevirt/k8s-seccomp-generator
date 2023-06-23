FROM quay.io/centos/centos:stream8
COPY ./installer.sh installer.sh
RUN curl -L -O https://download.falco.org/packages/bin/x86_64/falco-0.35.0-x86_64.tar.gz
RUN tar -xvf falco-0.35.0-x86_64.tar.gz
RUN mkdir /usr/_falco/ && mv falco-0.35.0-x86_64 /usr/_falco/falco
RUN chmod +x installer.sh
CMD installer.sh && cp /root/.falco/falco-bpf.o /usr/_falco/falco-bpf.o


