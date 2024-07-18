build:
	docker-compose build

start:
	docker-compose up ngrok api mongo -d
	docker-compose run cli

tests:
	docker-compose run tests go test ./...

stop:
	docker-compose down --remove-orphans