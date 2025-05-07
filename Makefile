.PHONY: build run test docker-up docker-down

build:
	go build -o bin/server ./cmd/server/main.go
	go build -o bin/worker ./cmd/worker/main.go

run: build
	./bin/server --db data.db --listen 8080 &
	./bin/worker --db data.db --poll 100ms

test:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

clean:
	rm -rf bin/ coverage.out coverage.html data.db