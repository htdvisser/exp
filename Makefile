TASK ?= tool/bin/task

.PHONY: default

default:

$(TASK): tool/task/main.go $(wildcard tool/task/commands/*.go) go.mod go.sum
	go build -o $@ $<

.PHONY: task

task: $(TASK)

.PHONY: deps.download

deps.download: | $(TASK)
	$(TASK) go mod download

.PHONY: deps.update

deps.update: | $(TASK)
	$(TASK) go get -u -t ./...
	$(TASK) go mod tidy

.PHONY: test

test: | $(TASK)
	$(TASK) go test ./...

.PHONY: cover

cover: | $(TASK)
	$(TASK) go test -covermode=atomic -coverprofile=coverage.out ./...
	$(TASK) go tool cover -html=coverage.out -o coverage.html

.PHONY: clean

clean:
	find . -name coverage.out -delete
	find . -name coverage.html -delete
