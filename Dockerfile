FROM alpine:latest
MAINTAINER andy@ak7.io

ADD bin/alexaops /bin/alexaops
ADD alexaops.conf.sample /etc/alexaops.conf

EXPOSE 80

CMD ["/bin/alexaops", "listen", "--address", ":80", "--config", "/etc/alexaops.conf"]