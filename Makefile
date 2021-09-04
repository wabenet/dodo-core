all: v1alpha1 v1alpha2 test lint

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

.PHONY: v1alpha1
v1alpha1: api/v1alpha1/plugin.pb.go api/v1alpha1/backdrop.pb.go api/v1alpha1/build.pb.go api/v1alpha1/configuration.pb.go api/v1alpha1/runtime.pb.go api/v1alpha1/builder.pb.go

.PHONY: v1alpha2
v1alpha2: api/v1alpha2/plugin.pb.go api/v1alpha2/backdrop.pb.go api/v1alpha2/build.pb.go api/v1alpha2/configuration.pb.go api/v1alpha2/runtime.pb.go api/v1alpha2/builder.pb.go

%.pb.go: %.proto
	protoc -I . --go_out=plugins=grpc:. --go_opt=module=github.com/dodo-cli/dodo-core $<
