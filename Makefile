.PHONY: test

run:
	docker-compose pull
	docker-compose up -d --build

test:
	docker-compose pull
	docker-compose up -d db
	go test ./... -p 1
