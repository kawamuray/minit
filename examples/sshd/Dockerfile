FROM kawamuray/minit
MAINTAINER Yuto Kawamura(kawamuray) <kawamuray.dadada@gmail.com>

RUN apt-get update
RUN apt-get install -y ssh
RUN sed -i 's/^\(# *\)\?PermitRootLogin .*$/PermitRootLogin yes/' /etc/ssh/sshd_config

RUN echo 'root:hoge' | chpasswd

ENTRYPOINT ["/minit", "--syslog", "--"]
CMD /etc/init.d/ssh start
