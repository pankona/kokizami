

.PHONY: all test build deps show

all: test

test: build
	@go test -coverprofile=coverage.out

build: deps
	@make -C $(CURDIR)/cmd/kkzm

deps: $(CURDIR)/vendor

$(CURDIR)/vendor:
	@glide install

glide:
ifeq ($(shell command -v glide 2> /dev/null),)
	go get -u github.com/Masterminds/glide
endif

install: test
	@make install -C $(CURDIR)/cmd/kkzm

show:
	@go tool cover -html=coverage.out
