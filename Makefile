build:
	CGO_ENABLED=0 go build -o ./bin/autoEqMac

lint:
	golangci-lint run
