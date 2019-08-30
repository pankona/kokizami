package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/pankona/kokizami"
	"github.com/pankona/kokizami/cmd/kkzm/repo"
	"github.com/pankona/kokizami/models"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = Name
	app.Version = Version
	app.Author = "pankona"
	app.Email = "yosuke.akatsuka@gmail.com"
	app.Usage = "awesome task timer and tracker"

	app.Flags = globalFlags()
	app.Commands = commands()
	app.Action = CmdList // show list if no argument
	app.CommandNotFound = CommandNotFound

	db := &sql.DB{}

	app.Before = func(ctx *cli.Context) error {
		u, err := user.Current()
		if err != nil {
			return fmt.Errorf("failed to get current user: %v", err)
		}

		configDir := filepath.Join(u.HomeDir, ".config", "kokizami")
		err = os.MkdirAll(configDir, 0755) // #nosec
		if err != nil {
			return fmt.Errorf("failed to create directory on %v", configDir)
		}

		db, err = openDB(filepath.Join(configDir, "db"))
		if err != nil {
			return fmt.Errorf("failed to open DB: %v", err)
		}

		err = repo.CreateTables(db)
		if err != nil {
			return fmt.Errorf("failed to create tables: %v", err)
		}

		kkzm := &kokizami.Kokizami{
			KizamiRepo:  repo.NewKizamiRepo(db),
			TagRepo:     repo.NewTagRepo(db),
			SummaryRepo: repo.NewSummaryRepo(db),
		}

		app.Metadata["kkzm"] = kkzm

		return nil
	}

	app.After = func(ctx *cli.Context) error {
		return db.Close()
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func openDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func enableVerboseQuery(enable bool) {
	models.XOLog = func(s string, p ...interface{}) {
		if enable {
			fmt.Printf("-------------------------------------\nQUERY: %s\n  VAL: %v\n", s, p)
		}
	}
}
