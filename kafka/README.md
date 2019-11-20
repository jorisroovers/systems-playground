
# Kafka
Kafka needs Zookeeper for data replication under the hood.

```sh
# Bring up Zookeeper and Kafka
docker-compose up
```

# Kafkacat
Simple program to interact with kafka via CLI

```sh
# Consuming
alias kafkacat='docker run -ti --network=$(docker network ls -f name=kafka -q) edenhill/kafkacat:1.5.0'
kafkacat -b kafka -L
kafkacat -b kafka -t my-topic

# Producing (=separate terminal)
# IMPORTANT: slightly different kafkacat alias: no passing of -t, so we can pipe echo commands into the container
# This fixes the "the input device is not a TTY" error message.
alias kafkacat='docker run -i --network=$(docker network ls -f name=kafka -q) edenhill/kafkacat:1.5.0'
echo "foobar" | kafkacat -b kafka -t my-topic -p 0
```

# Python example

```sh
docker build -t kafka-python-example -f Dockerfile.python .

# Container alias: connect to kafka network, mount current dir under /data
alias container='docker run --network=$(docker network ls -f name=kafka -q) -ti -v $(pwd):/data kafka-python-example'
container python /data/kafka.py
```