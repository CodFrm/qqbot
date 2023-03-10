
arm:
	GOOS=linux GOARCH=arm go build -o qqbot ./cmd/app

linux:
	GOOS=linux go build -o qqbot ./cmd/app

migrate:
	migrate create -ext sql -dir ./db/migrations $(MIGRATE)
