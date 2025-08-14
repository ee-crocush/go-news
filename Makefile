.PHONY: build up down restart status restart-consumers

build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down

restart: down up

status:
	docker compose ps

restart-consumers:
	docker compose restart news-comments
	docker compose restart news-moderation