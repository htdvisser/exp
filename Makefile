default:

modules = $(patsubst %/go.mod,%,$(shell find . -mindepth 2 -name "go.mod" | sort))

%/coverage.out:
	cd $(shell dirname $@); go test -covermode=atomic -coverprofile=coverage.out ./...

coverfiles = $(patsubst %,%/coverage.out,$(modules))

coverage.out: $(coverfiles)
	echo "mode: atomic" > $@
	tail -qn+2 $(coverfiles) >> $@

coverage.html: coverage.out
	go tool cover -html=$< -o $@

clean:
	rm -f $(coverfiles)
	rm -f coverage.out
	rm -f coverage.html
