FROM alpine:latest
MAINTAINER andy@ak7.io

ADD bin/alexaops /bin/alexaops

EXPOSE 80

CMD ["/bin/alexaops", "listen", "--address", ":80"]