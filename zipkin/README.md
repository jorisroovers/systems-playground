# Zipkin

To run zipkin:
```bash
docker-compose up
```
The ```docker-compose.yml``` file contains 2 containers (linked with a network):
1. Zipkin, hosted on http://localhost:9411 on the host
2. MySQL, which zipkin uses as transport layer. You can also use Kafka, Cassandra, Elasticsearch and others as transport layers.


# Zipkin transport

```bash
docker run -d -p 2181:2181 openzipkin/zipkin-kafka

docker run -d -p 9411:9411 openzipkin/zipkin

virtualenv .venv && source .venv/bin/activate
pip install -r requirements.txt
```