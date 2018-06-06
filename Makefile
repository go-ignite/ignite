run:
	go run ./cmd/ignite/main.go

docs:
	swag init ./handler -g ./cmd/ignite/main.go
