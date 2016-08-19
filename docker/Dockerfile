FROM alpine:3.2

MAINTAINER Otto Giron <ottog2486@gmail.com>
RUN apk update
ADD build/dist/linux/ferraristream /
ENTRYPOINT ["/ferraristream"]
