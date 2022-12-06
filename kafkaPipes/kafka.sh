bin/zookeeper-server-start.sh config/zookeeper.properties
# New terminal
bin/kafka-server-start.sh config/server.properties
# New terminal
bin/kafka-topics.sh --create --topic quickstart-events --bootstrap-server localhost:9092
bin/kafka-console-consumer.sh --topic quickstart-events --from-beginning --bootstrap-server localhost:9092
