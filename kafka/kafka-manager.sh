#!/bin/bash

# Цвета для красоты
GREEN='\033[0;32m'
NC='\033[0m'
DELIMITER='------------------------------------------'
CONTAINTER=news-kafka

function show_menu() {
  echo -e "${GREEN}Выберите действие:${NC}"
  echo "1) Создать топик"
  echo "2) Посмотреть все топики"
  echo "3) Отправить сообщение (producer)"
  echo "4) Прочитать сообщения (consumer)"
  echo "5) Войти в контейнер Kafka"
  echo "6) Выйти"
  echo
}

function create_topic() {
  read -p "Название топика (по умолчанию: test-topic): " TOPIC
  TOPIC=${TOPIC:-test-topic}
  read -p "Количество партиций (по умолчанию: 1): " PARTITIONS
  PARTITIONS=${PARTITIONS:-1}
  read -p "Фактор репликации (по умолчанию: 1): " REPLICATION
  REPLICATION=${REPLICATION:-1}

  echo -e "${GREEN}Создание топика '${TOPIC}'...${NC}"
  docker exec ${CONTAINTER} /usr/bin/kafka-topics \
    --create \
    --topic "$TOPIC" \
    --bootstrap-server localhost:9092 \
    --partitions "$PARTITIONS" \
    --replication-factor "$REPLICATION"
}

function list_topics() {
  echo -e "${GREEN}Список топиков:${NC}"
  docker exec ${CONTAINTER} /usr/bin/kafka-topics --list --bootstrap-server localhost:9092
  echo ${DELIMITER}
}

function produce_messages() {
  read -p "Название топика (по умолчанию: test-topic): " TOPIC
  TOPIC=${TOPIC:-test-topic}

  echo -e "${GREEN}Пишите сообщения. Для выхода нажмите Ctrl+C.${NC}"
  docker exec -it kafka /usr/bin/kafka-console-producer \
    --broker-list localhost:9092 \
    --topic "$TOPIC"
}

function consume_messages() {
  read -p "Название топика (по умолчанию: test-topic): " TOPIC
  TOPIC=${TOPIC:-test-topic}

  echo -e "${GREEN}Чтение сообщений из '${TOPIC}'... Нажмите Ctrl+C для выхода.${NC}"
  docker exec -it ${CONTAINTER} /usr/bin/kafka-console-consumer \
    --bootstrap-server localhost:9092 \
    --topic "$TOPIC" \
    --from-beginning
}

function open_shell() {
  docker exec -it ${CONTAINTER} bash
}

# Основной цикл
while true; do
  show_menu
  read -p "Ваш выбор (1-6): " CHOICE

  case $CHOICE in
    1) create_topic ;;
    2) list_topics ;;
    3) produce_messages ;;
    4) consume_messages ;;
    5) open_shell ;;
    6) echo "Выход..."; break ;;
    *) echo "Неверный выбор. Попробуйте снова." ;;
  esac

  echo
done
