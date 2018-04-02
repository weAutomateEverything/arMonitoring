FROM alpine:3.6
RUN apk update
RUN apk add cifs-utils
RUN mkdir /mnt/share01 && mkdir /mnt/share02