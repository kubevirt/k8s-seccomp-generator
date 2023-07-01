FROM golang:1.20
RUN cd /; git clone https://github.com/sudo-NithishKarthik/k8s_seccomp_generator.git
WORKDIR /k8s_seccomp_generator
RUN cd seccomp_gen/cmd/falco_syscalls_formatter && go build -o formatter
RUN cd seccomp_gen/cmd/syscalls_tracer && go build -o tracer
CMD cp seccomp_gen/cmd/falco_syscalls_formatter/formatter /falco/formatter && seccomp_gen/cmd/syscalls_tracer/tracer
