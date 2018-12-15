DEPS = $(wildcard */*.go)
VERSION = $(shell git describe --always --dirty)

all: imports lint vet posd

posd: main.go $(DEPS)
	CGO_ENABLED=0 GOOS=linux \
	  go build -a \
		  -ldflags="-X main.posdVersion=$(VERSION)" \
	    -installsuffix cgo -o $@ $<
	strip $@

clean:
	rm -f posd

lint:
	@ go get -v github.com/golang/lint/golint
	@for file in $$(git ls-files '*.go' | grep -v '_workspace/' | grep -v 'vendor/'); do \
		export output="$$(golint $${file} | grep -v 'type name will be used as docker.DockerInfo')"; \
		[ -n "$${output}" ] && echo "$${output}" && export status=1; \
	done; \
	exit $${status:-0}

vet: main.go
	go vet $<

imports: main.go
	dep ensure -vendor-only
	goimports -d $<

.PHONY: all lint clean
