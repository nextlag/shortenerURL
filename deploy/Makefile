build:
	docker build -t short ..
run:
	docker run -p 8080:8080 short
test:
	docker-compose -f ./Docker-Compose.yaml up --build --abort-on-container-exit
	docker-compose -f ./Docker-Compose.yaml down --volumes