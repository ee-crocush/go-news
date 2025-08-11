#!/bin/bash

# Цвета для красоты
GREEN='\033[0;32m'
NC='\033[0m'
DELIMITER='------------------------------------------'

# Параметры по умолчанию
PARTITIONS=1
REPLICATION=1

function create_topic() {
  local TOPIC="$1"
  echo -e "${GREEN}Создание топика '${TOPIC}'...${NC}"
  docker exec kafka /usr/bin/kafka-topics \
    --create \
    --topic "$TOPIC" \
    --bootstrap-server localhost:9092 \
    --partitions "$PARTITIONS" \
    --replication-factor "$REPLICATION"
  echo ${DELIMITER}
}

# Создаём два топика
create_topic "comments.created"
create_topic "comments.moderated"

echo -e "${GREEN}Все топики созданы успешно.${NC}"