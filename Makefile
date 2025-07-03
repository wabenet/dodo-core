.PHONY: all
all: clean test build lint

.PHONY: clean
clean:
	rm -rf ./dist c.out coverage.txt

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: copier-update
copier-update:
	copier update --trust --skip-answered

.PHONY: update
update:
	go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all | xargs --no-run-if-empty go get
	go mod tidy

.PHONY: build
build: build-proto

.PHONY: release
release:

.PHONY: test
test: test-go

.PHONY: lint
lint: lint-proto lint-go

.PHONY: coverage-report
coverage-report:
	cc-test-reporter before-build
	go test -race -coverprofile=coverage.txt -covermode=atomic -coverpkg=./... ./...
	cp coverage.txt c.out
	cc-test-reporter after-build -t gocov -p $$(go list -m)

.PHONY: test-go
test-go: build-proto
	go test -race -cover ./...

.PHONY: lint-go
lint-go:
	golangci-lint run

.PHONY: lint-proto
lint-proto:
	buf lint

.PHONY: build-proto
build-proto:
	buf generate
