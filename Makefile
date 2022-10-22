VERSION=$(shell cat VERSION)
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_HASH=$(shell git rev-parse HEAD)
CONTAINER_PLATFORM?=linux/amd64

BUILDS=linux.amd64 linux.386 linux.arm64 linux.mips64 \
	   windows.amd64.exe \
	   freebsd.amd64 \
	   darwin.amd64 darwin.arm64

BINARIES=$(addprefix bin/r2t-$(VERSION)., $(BUILDS))

LDFLAGS="-X main.Version=$(VERSION) -X main.BuildDate=$(BUILD_DATE) -X main.GitHash=$(GIT_HASH)"

toc:
	@echo "list of targets:"
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | \
		awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | \
		sort | \
		egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | \
		awk '{ print " ", $$1 }'

r2t: bin/r2t

bin/r2t: bin .
	go build -ldflags $(LDFLAGS) -v -o $@ .

######################################################
## release related

release: $(BINARIES)

bin/r2t-$(VERSION).linux.%: bin
	env GOOS=linux GOARCH=$* CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -o $@ .

bin/r2t-$(VERSION).darwin.%: bin
	env GOOS=darwin GOARCH=$* CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -o $@ .

bin/r2t-$(VERSION).windows.%.exe: bin
	env GOOS=windows GOARCH=$* CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -o $@ .

bin/r2t-$(VERSION).freebsd.%: bin
	env GOOS=freebsd GOARCH=$* CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -o $@ .

bin:
	mkdir $@

container-image:
	env DOCKER_BUILDKIT=1 docker build \
		--file Dockerfile \
		--platform=$(CONTAINER_PLATFORM) \
		--build-arg VERSION=$(VERSION) \
		--tag $(CONTAINER_PLATFORM)-r2t:$(VERSION) .

######################################################
## dev related

compile-analysis: .
	go build -gcflags '-m' ./$^

code-quality:
	@echo "--"
	-go vet .
	@echo "--"
	-staticcheck .
	@echo "--"
	-gofmt -s -d .
	@echo "--"
	-gocyclo .
	@echo "--"
	-ineffassign .

fetch-report-tools:
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install github.com/client9/misspell/cmd/misspell@latest
	go install github.com/gordonklaus/ineffassign@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

test:
	go test -v .
