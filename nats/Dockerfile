FROM alpine:3.10.3
WORKDIR /tmp

RUN apk add curl

RUN curl -LOs https://github.com/nats-io/go-nats-examples/releases/download/0.1.0/go-nats-examples-v0.1.0-linux-amd64.zip
RUN unzip go-nats-examples-v0.1.0-linux-amd64.zip

# NATS requires glibc, alpine uses musl. Symlink musl .so to glib .so path. This works because libc API is standardized.
# Details: https://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

RUN cp go-nats-examples-v0.1.0-linux-amd64/nats-* /usr/bin/

