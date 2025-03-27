SHELL:=/bin/bash
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

# include .env.test

build:
	cd $(ROOT_DIR)/httpserver && \
	go build -o brick-monitoring ./cmd/main.go

dep-dev:
	docker compose --env-file .env.dev -p dev --file docker-compose-dependencies.yml up --build -d --force-recreate

server-dev: build
	./httpserver/brick-monitoring --env-file=.env.dev

# Запуск в dev окружении
dev:
	make dep-dev
# TODO: fix Hardcode
	sleep 1
	make server-dev

# запуск тестового окружения для интеграционных тестов
dep-test:
	docker compose --env-file .env.test -p test --file docker-compose-dependencies.yml up --build --force-recreate

