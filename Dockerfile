FROM golang
MAINTAINER Yuto Kawamura(kawamuray) <kawamuray.dadada@gmail.com>

# Install minit
RUN git clone https://github.com/kawamuray/minit.git /tmp/minit
RUN make -C /tmp/minit
RUN mv /tmp/minit/minit /minit
