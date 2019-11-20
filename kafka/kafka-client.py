from kafka import KafkaConsumer

topic = "my-topic"
consumer = KafkaConsumer(topic, bootstrap_servers="kafka")
print(f"Listening for message on topic '{topic}'")
for msg in consumer:
    print (msg)

# Producing messages is very similar using KafkaProducer