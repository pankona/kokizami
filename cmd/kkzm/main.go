package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/pankona/kokizami"
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

	kkzm := &kokizami.Kokizami{}

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

		db, err := openDB(filepath.Join(configDir, "db"))
		if err != nil {
			return fmt.Errorf("failed to open DB: %v", err)
		}
		err = createTables(db)
		if err != nil {
			return fmt.Errorf("failed to create tables: %v", err)
		}

		err = kkzm.Initialize()
		if err != nil {
			return fmt.Errorf("failed to initialize kkzm: %v", err)
		}
		kkzm.SetDB(db)
		kkzm.KizamiRepo = &kizamiRepo{
			db:  db,
			now: time.Now,
		}
		kkzm.TagRepo = &tagRepo{db: db}
		kkzm.SummaryRepo = &summaryRepo{db: db}

		app.Metadata["kkzm"] = kkzm
		return nil
	}

	app.After = func(ctx *cli.Context) error {
		return kkzm.Finalize()
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

func createTables(db *sql.DB) error {
	if err := models.CreateKizamiTable(db); err != nil {
		return fmt.Errorf("failed to create kizami table: %v", err)
	}
	if err := models.CreateTagTable(db); err != nil {
		return fmt.Errorf("failed to create tag table: %v", err)
	}
	if err := models.CreateRelationTable(db); err != nil {
		return fmt.Errorf("failed to create relation table: %v", err)
	}
	return nil
}

// EnableVerboseQuery toggles debug logging by argument
func EnableVerboseQuery(enable bool) {
	models.XOLog = func(s string, p ...interface{}) {
		if enable {
			fmt.Printf("-------------------------------------\nQUERY: %s\n  VAL: %v\n", s, p)
		}
	}
}
