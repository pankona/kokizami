# kokizami

## Description

this application is for tracking tasks and its elapsed time.

## Usage

following commands are available.

* kkzm start [desc]
    * if desc is NOT specified, enter mode to edit desc via editor 
* kkzm stop [kizami id]
* kkzm restart (kizami id)
* kkzm list
* kkzm edit [id] [desc|started_at|stopped_at] [new value]
    * if desc (or something) is NOT specified, enter mode to edit desc, started_at and stopped_at via editor 

## Install

To install, use `go get`:

```bash
$ go get -u github.com/pankona/kokizami/cmd/kkzm
```

## Contribution

1. Fork ([https://github.com/pankona/kokizami/fork](https://github.com/pankona/kokizami/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## License

MIT

## Author

[pankona](https://github.com/pankona)
