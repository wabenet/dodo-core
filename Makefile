all: proto test lint

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
	go test -cover ./...

.PHONY: proto
proto: api/v1alpha1/plugin.pb.go api/v1alpha1/backdrop.pb.go api/v1alpha1/build.pb.go api/v1alpha1/configuration.pb.go api/v1alpha1/runtime.pb.go api/v1alpha1/builder.pb.go

%.pb.go: %.proto
	protoc -I . --go_out=plugins=grpc:. --go_opt=module=github.com/dodo-cli/dodo-core $<
