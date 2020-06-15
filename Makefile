
arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build .

migrate:
	migrate create -ext sql -dir ./db/migrations $(MIGRATE)
