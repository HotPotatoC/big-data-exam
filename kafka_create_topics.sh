docker exec broker \
    kafka-topics --create \
    --topic messages \
    --bootstrap-server broker:9092 \
    --partitions 1 \
    --replication-factor 1

docker exec broker \
    kafka-topics --create \
    --topic transactions \
    --bootstrap-server broker:9092 \
    --partitions 1 \
    --replication-factor 1
