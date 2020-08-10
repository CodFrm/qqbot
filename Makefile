
arm:
	GOOS=linux GOARCH=arm go build .

migrate:
	migrate create -ext sql -dir ./db/migrations $(MIGRATE)
