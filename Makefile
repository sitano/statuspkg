PKGS := github.com/sitano/statuspkg
SRCDIRS := $(shell go list -f '{{.Dir}}' $(PKGS))
GO := go

check: test vet gofmt unconvert staticcheck ineffassign unparam

test:
	env GO111MODULE=on $(GO) test $(PKGS)

vet: | test
	$(GO) vet $(PKGS)

staticcheck:
	env GO111MODULE=off $(GO) get honnef.co/go/tools/cmd/staticcheck
	staticcheck -checks all $(PKGS)

misspell:
	env GO111MODULE=off $(GO) get github.com/client9/misspell/cmd/misspell
	misspell \
		-locale GB \
		-error \
		*.md *.go

unconvert:
	env GO111MODULE=off $(GO) get github.com/mdempsky/unconvert
	unconvert -v $(PKGS)

ineffassign:
	env GO111MODULE=off $(GO) get github.com/gordonklaus/ineffassign
	find $(SRCDIRS) -name '*.go' | xargs ineffassign

pedantic: check errcheck

unparam:
	env GO111MODULE=off $(GO) get mvdan.cc/unparam
	unparam ./...

errcheck:
	env GO111MODULE=off $(GO) get github.com/kisielk/errcheck
	errcheck $(PKGS)

gofmt:
	@echo Checking code is gofmted
	@test -z "$(shell gofmt -s -l -d -e $(SRCDIRS) | tee /dev/stderr)"
