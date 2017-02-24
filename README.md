# kokizami

## Description

this application is for tracking tasks (that is called *kizami*) and its elapsed time (like toggl).
since this application is standalone, network access is not required.

## Usage

following commands are available.

* kkzm start [desc]
    * if desc is NOT specified, enter mode to edit desc via editor.
* kkzm stop [kizami id]
    * if id is not specified, all kizami that is in ongoing are stopped.
* kkzm restart (kizami id)
    * start new kizami with specified id's description.
* kkzm list
    * lists all kizamis
* kkzm edit [id] [desc|started_at|stopped_at] [new value]
    * if desc (or something) is NOT specified, open editor to edit all of them.
* kkzm delete [id]
    * delete specified kizami

## Notes

* this application will create a database file (.kokizami.db) at user's home directory.

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
