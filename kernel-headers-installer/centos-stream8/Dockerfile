FROM quay.io/centos/centos:stream8
RUN curl -L -O https://download.falco.org/packages/bin/x86_64/falco-0.35.0-x86_64.tar.gz
RUN tar -xvf falco-0.35.0-x86_64.tar.gz
COPY ./installer.sh /root/installer.sh
RUN chmod 777 /root/installer.sh
CMD /root/installer.sh && mkdir /usr/_falco-probe/ && cp /root/.falco/falco-bpf.o /usr/_falco-probe/falco-bpf.o


