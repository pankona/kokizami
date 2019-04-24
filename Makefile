

all: build test lint

build:
	@make -C $(CURDIR)/cmd/kkzm

lint:
	gometalinter --vendor --exclude="xo.go" ./...

test: install-goverage
	@go test -cover ./...
	@goverage -coverprofile=coverage.out ./...

install:
	@make install -C $(CURDIR)/cmd/kkzm

show-coverage:
	@go tool cover -html=coverage.out

install-goverage:
ifeq ($(shell command -v goverage 2> /dev/null),)
	go get -u github.com/haya14busa/goverage
endif
