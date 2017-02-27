

PKGS=$(shell go list ./...  | grep -v vendor)

.PHONY: all test build deps show

all: test

test: build
	@golint $(PKGS)
	@gosimple $(PKGS)
	@go test -cover $(PKGS)
	@goverage -coverprofile=coverage.out $(PKGS)

build: deps
	@make -C $(CURDIR)/cmd/kkzm

deps: $(CURDIR)/vendor

$(CURDIR)/vendor:
	@glide install

install-tools: glide goverage
	@echo done!

glide:
ifeq ($(shell command -v glide 2> /dev/null),)
	go get -u github.com/Masterminds/glide
endif

goverage:
ifeq ($(shell command -v goverage 2> /dev/null),)
	go get -u github.com/haya14busa/goverage
endif

install: test
	@make install -C $(CURDIR)/cmd/kkzm

show:
	@go tool cover -html=coverage.out
