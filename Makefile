

.PHONY: all test build show

all: test

test: build
	@go test -coverprofile=coverage.out

build:
	@make -C $(CURDIR)/cmd/todo

install:
	@make install -C $(CURDIR)/cmd/todo

show:
	@go tool cover -html=coverage.out
