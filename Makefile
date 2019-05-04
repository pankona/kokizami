
COVFILE=.coverage.out

all: build test lint

build:
	@make -C $(CURDIR)/cmd/kkzm

lint:
	gometalinter --vendor --exclude="xo.go" ./...

test: install-goverage
	@go test -cover ./...
	@goverage -coverprofile=$(COVFILE) ./...

install:
	@make install -C $(CURDIR)/cmd/kkzm

show-coverage:
	@go tool cover -html=$(COVFILE)

install-goverage:
ifeq ($(shell command -v goverage 2> /dev/null),)
	go get -u github.com/haya14busa/goverage
endif
