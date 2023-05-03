run:
	docker-compose pull
	docker-compose up -d --build

test: 
	docker-compose pull
	docker-compose up -d db
	cd api; go test ./... -p 1
