shorten:
	docker run --name shorten -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=Xer_0101 -d postgres:12-alpine
	sleep 10
	$(MAKE) db_create

db_create:
	docker exec -it shorten psql -U postgres -c "CREATE DATABASE shorten;"

run:
	go run cmd/shortener/main.go -d "postgres://postgres:Xer_0101@localhost/shorten?sslmode=disable"

runF:
	go run cmd/shortener/main.go -f file.json
build:
	docker build -t short ..
docker_run:
	docker run -p 8080:8080 short
test:
	docker-compose -f ./Docker-Compose.yaml up --build --abort-on-container-exit
	docker-compose -f ./Docker-Compose.yaml down --volumes
lint:
	go run ./cmd/staticlint ./...

.PHONY: shorten db_create run runF build docker_run test lint