PROJECT_NAME=country-app
SHELL=/bin/zsh

setup:
	docker-compose -p $(PROJECT_NAME) up -d

destroy:
	docker-compose down --remove-orphans

run:
	go run main.go db.go es.go server.go
