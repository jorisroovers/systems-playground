# NATS
Cloud-Native Messaging Queue. Client libraries available for Go, Ruby, Java, Python, and more.

## Server
```sh
docker pull nats:2.1.0-alpine3.10
docker run -p 4222:4222 -ti nats:2.1.0-alpine3.10
```

## Client
```sh
# Build client docker image, based on Alpine
# Make sure you're in the NATS directory
docker build -t nats-cli .

# Start listener on all topics (= ">")
# Note the use of docker 'host' networking to be able to connect to the NATS server
docker run --net=host nats-cli nats-sub -s nats://localhost:4222 ">"

# Publish Hello World Message on 'my-topic'
docker run --net=host nats-cli nats-pub -s nats://localhost:4222 my-topic "Hello World"

# Only listen to messages on "my-topic"
docker run --net=host nats-cli nats-sub -s nats://localhost:4222 "my-topic"

# Hierarchical message structure and wildcard matching:
docker run --net=host nats-cli nats-sub -s nats://localhost:4222 "my-topic.bar"
docker run --net=host nats-cli nats-sub -s nats://localhost:4222 "my-topic.*"
docker run --net=host nats-cli nats-sub -s nats://localhost:4222 "my-topic.*.foo"
```