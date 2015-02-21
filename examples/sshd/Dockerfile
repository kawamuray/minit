FROM centos:centos6
MAINTAINER Yuto Kawamura(kawamuray) <kawamuray.dadada@gmail.com>

RUN yum install -y tar openssh-server openssh-clients git

RUN echo hoge | passwd --stdin root

# Install Go
RUN curl -sL https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz | tar zx -C /usr/local
ENV PATH /usr/local/go/bin:$PATH

# Install minit
RUN git clone https://github.com/kawamuray/minit.git /tmp/minit
RUN make -C /tmp/minit
RUN mv /tmp/minit/minit /minit

CMD /minit --syslog /etc/init.d/sshd start
