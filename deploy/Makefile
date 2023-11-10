build:
	docker build -t short ..
run:
	docker run -p 8080:8080 short
test:
	docker-compose -f ./Docker-Compose.yaml up --build --abort-on-container-exit
	docker-compose -f ./Docker-Compose.yaml down --volumes
postgres:
	docker run --name shorten -p 5432:5432 -e POSTGRES_USER=shorten -e POSTGRES_PASSWORD=skypass12345 -d postgres:12-alpine