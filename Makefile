SHELL:=/bin/bash
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

# include .env.test

dependencies-dev:
	docker compose --env-file .env.dev -p dev --file dependencies/docker-compose.yml up --build -d --force-recreate

app-dev:
	docker compose --env-file .env.dev -p dev --file docker-compose.yml up --build -d --force-recreate
# server-build:
# 	cd $(ROOT_DIR)/httpserver && \
# 	go build -o brick-monitoring ./cmd/main.go

# server-dev: build
# 	./httpserver/brick-monitoring --env-file=.env.dev

# producer-dev:
# 	. $(ROOT_DIR)/producer/.venv/bin/activate && \
# 	python3 $(ROOT_DIR)/producer/producer.py

# translator-build:
# 	cd $(ROOT_DIR)/translator && \
# 	go build -o translator ./app.go

# translator-dev: translator-build
# 	./translator/translator

# Запуск в dev окружении
dev:
	docker compose --env-file .env.dev -p dev --file dependencies/docker-compose.yml up --build -d --force-recreate
	sleep 5
	docker compose --env-file .env.dev -p dev --file docker-compose.yml up --build -d --force-recreate

# запуск тестового окружения для интеграционных тестов
# dependencies-test:
# 	docker compose --env-file .env.test -p test --file docker-compose-dependencies.yml up --build --force-recreate
