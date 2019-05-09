run:
	source .env && go run ./cmd/ignite/main.go

docs:
	swag init ./handler -g ./cmd/ignite/main.go

link-agent:
	rm -rf vendor/github.com/go-ignite/ignite-agent
	ln -s ${GOPATH}/src/github.com/go-ignite/ignite-agent vendor/github.com/go-ignite
