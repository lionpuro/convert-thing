build:
	docker compose build

run:
	docker compose up --build

fmt:
	@gofmt -l -s -w .
