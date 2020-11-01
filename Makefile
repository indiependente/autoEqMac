.PHONY: build
build:
	CGO_ENABLED=0 go build -o ./bin/autoEqMac

.PHONY: coverage
coverage:
	gopherbadger -md="README.md" -png=false

.PHONY: lint
lint:
	golangci-lint run
