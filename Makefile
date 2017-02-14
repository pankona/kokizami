

.PHONY: all test build show

all: test

test: build
	@go test -coverprofile=coverage.out

build:
	@make -C $(CURDIR)/cmd/kkzm

install:
	@make install -C $(CURDIR)/cmd/kkzm

show:
	@go tool cover -html=coverage.out
