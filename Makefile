PWD = ${CURDIR}
# Название сервиса
SERVICE_NAME = cathedral-bot


# компиляция сервиса
.PHONY: build
build:
	go build -o bin/$(SERVICE_NAME)  $(PWD)/cmd/$(SERVICE_NAME)  -config ./configs/config.yaml

# Запуск сервиса
.PHONY: run
run:
	go run $(PWD)/cmd/$(SERVICE_NAME) -config ./configs/config.yaml

# Запуск миграций
.PHONY: migrate
migrate:
	go run $(PWD)/cmd/migrate -config=./configs/config.yaml
