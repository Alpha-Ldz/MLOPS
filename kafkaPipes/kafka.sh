systemctl start docker
docker-compose up -d
docker exec broker kafka-topics --bootstrap-server broker:9092 --create --topic rapport-topic
