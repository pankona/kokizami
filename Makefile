
COVFILE=$(CURDIR)/.coverage.out

all: build test lint

build:
	@make -C $(CURDIR)/cmd/kkzm

lint:
	golangci-lint run --new-from-rev= --deadline 300s

test: install-goverage
	@go test -cover ./...
	@goverage -coverprofile=$(COVFILE) ./...

install:
	@make install -C $(CURDIR)/cmd/kkzm

show-coverage: $(COVFILE)
	@go tool cover -html=$(COVFILE)

$(COVFILE):
	make test

install-goverage:
ifeq ($(shell command -v goverage 2> /dev/null),)
	go get -u github.com/haya14busa/goverage
endif

clean:
	rm -f $(COVFILE)
