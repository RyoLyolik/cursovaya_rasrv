SHELL:=/bin/bash
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

# include .env.test

# Запуск в dev окружении
dev:
	docker compose --env-file .env.dev -p dev --file dependencies/docker-compose.yml up --build -d --force-recreate
	sleep 5
	docker compose --env-file .env.dev -p dev --file docker-compose.yml up --build -d --force-recreate

# запуск тестового окружения для интеграционных тестов
# dependencies-test:
# 	docker compose --env-file .env.test -p test --file docker-compose-dependencies.yml up --build --force-recreate
