#!/bin/sh

set -e

# Параметры топиков
TOPICS="comments.created comments.moderated"
PARTITIONS=1
REPLICATION=1

echo "Ждём Kafka..."

while ! nc -z localhost 9092; do
  echo "Ожидаем Kafka на localhost:9092..."
  sleep 2
done

echo "Kafka запущена, создаём топики..."

for TOPIC in $TOPICS; do
  kafka-topics --create \
    --bootstrap-server localhost:9092 \
    --replication-factor $REPLICATION \
    --partitions $PARTITIONS \
    --topic "$TOPIC" || echo "Топик $TOPIC уже существует"
done

echo "Топики созданы, запускаем Kafka..."

exec /etc/confluent/docker/run