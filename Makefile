
arm:
	GOOS=linux GOARCH=arm go build .

linux:
	GOOS=linux go build .

migrate:
	migrate create -ext sql -dir ./db/migrations $(MIGRATE)
