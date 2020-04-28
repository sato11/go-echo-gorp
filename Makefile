.PHONY: deps
deps:
	go mod download

.PHONY: devel-deps
devel-deps:
	go get golang.org/x/lint/golint

.PHONY: lint
lint:
	go vet ./...
	golint -set_exit_status ./...

.PHONY: build
build:
	go build

.PHONY: run
run: build
	./go-echo-gorp
