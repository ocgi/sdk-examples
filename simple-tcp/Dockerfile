FROM centos:centos7

# install some tools for debug
RUN yum install -y net-tools telnet iputils

COPY ./bin/simple-tcp /simple-tcp

ENTRYPOINT ["/simple-tcp"]
