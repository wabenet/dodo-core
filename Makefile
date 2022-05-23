all: v1alpha3 test lint

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -cover -race ./pkg/...

.PHONY: v1alpha3
v1alpha3: api/v1alpha3/plugin.pb.go api/v1alpha3/backdrop.pb.go api/v1alpha3/build.pb.go api/v1alpha3/configuration.pb.go api/v1alpha3/runtime.pb.go api/v1alpha3/builder.pb.go

%.pb.go: %.proto
	protoc -I . --go_out=plugins=grpc:. --go_opt=module=github.com/dodo-cli/dodo-core $<
