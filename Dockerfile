FROM alpine:latest
MAINTAINER andy@ak7.io

RUN mkdir /var/alexaops
WORKDIR /var/alexaops

ADD bin/alexaops /bin/alexaops
ADD alexaops.conf.sample /var/alexaops/alexaops.conf

VOLUME /var/alexaops

EXPOSE 80

ENTRYPOINT ["/bin/alexaops"]