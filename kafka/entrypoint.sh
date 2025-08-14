#!/bin/sh
set -e

TOPICS="comments.created comments.moderated"
PARTITIONS=1
REPLICATION=1
KAFKA_BROKER=${KAFKA_BROKER:-news-kafka:9092}



exec /etc/confluent/docker/run

echo "Ждём Kafka на $KAFKA_BROKER..."

# Ждём, пока KRaft контроллер и брокер не станут готовыми
while true; do
  # проверяем, есть ли хотя бы один node
  if kafka-topics --bootstrap-server "$KAFKA_BROKER" --list >/dev/null 2>&1; then
    break
  fi
  echo "Ожидаем регистрацию брокера..."
  sleep 2
done

echo "Kafka готова, создаём топики..."

for TOPIC in $TOPICS; do
  kafka-topics --create \
    --bootstrap-server "$KAFKA_BROKER" \
    --replication-factor $REPLICATION \
    --partitions $PARTITIONS \
    --topic "$TOPIC" || echo "Топик $TOPIC уже существует"
done

echo "Топики созданы, запускаем Kafka..."