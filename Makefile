GOLANG_CI_VERSION ?= 'v1.64.4'

build:
	go build -a -installsuffix cgo -o ipt-netflow-exporter ./cmd/iptnetflowexporter

mocks:
	go run github.com/vektra/mockery/v2@v2.52.3

tests:
	go test -v ./...

lint:
	golangci-lint run -v ./...

install_linter:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin $(GOLANG_CI_VERSION)

