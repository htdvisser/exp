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
	go get -u -t ./...
	$(TASK) go mod tidy -go=1.16
	$(TASK) go mod tidy -go=1.17
	go mod tidy
	go work sync

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

.github/dependabot.yml: tool/dependabot.yml.tmpl $(shell find . -name go.mod) | $(TASK)
	$(TASK) gen dependabot

go.work: $(shell find . -name go.mod) | $(TASK)
	$(TASK) gen go-workspace
