# kokizami

## Description

This application is for tracking tasks (each task is called _kizami_ in this application) and its elapsed time.
Since this application is standalone, network access is not required.

## Usage

Following commands are available.

```
NAME:
   kkzm - awesome task timer and tracker

USAGE:
   kkzm [global options] command [command options] [arguments...]

VERSION:
   2.0.0

AUTHOR:
   pankona <yosuke.akatsuka@gmail.com>

COMMANDS:
     start    Start new task
     restart  Restart old task
     edit     Edit task
     list     Show list of tasks
     stop     Stop task
     delete   Delete task
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Notes

- This application will create a database file on `$HOME/.config/kokizami/db`

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

[Yosuke Akatsuka (a.k.a pankona)](https://github.com/pankona)
