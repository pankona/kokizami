

.PHONY: test
test:
	go test -coverprofile=coverage.out


.PHONY: show
show:
	go tool cover -html=coverage.out
