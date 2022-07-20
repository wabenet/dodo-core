all: v1alpha4 test lint

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

.PHONY: v1alpha4
v1alpha4: api/v1alpha4/plugin.pb.go api/v1alpha4/backdrop.pb.go api/v1alpha4/build.pb.go api/v1alpha4/configuration.pb.go api/v1alpha4/runtime.pb.go api/v1alpha4/builder.pb.go

%.pb.go: %.proto
	protoc -I . --go_out=plugins=grpc:. --go_opt=module=github.com/wabenet/dodo-core $<
