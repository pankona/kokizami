package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/pankona/kokizami"
	"github.com/urfave/cli"
)

// GlobalFlags can be used globally
var GlobalFlags = []cli.Flag{}

// Commands represents list of subcommands
var Commands = []cli.Command{
	{
		Name:   "start",
		Usage:  "start new task",
		Action: CmdStart,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "restart",
		Usage:  "restart old task",
		Action: CmdRestart,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "edit",
		Usage:  "edit task",
		Action: CmdEdit,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "list",
		Usage:  "show list of tasks",
		Action: CmdList,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "stop",
		Usage:  "stop task",
		Action: CmdStop,
		Flags:  []cli.Flag{},
	},
	{
		Name:   "delete",
		Usage:  "delete task",
		Action: CmdDelete,
		Flags:  []cli.Flag{},
	},
}

// CommandNotFound is called when specified subcommand is not found
func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}

func round(d, r time.Duration) time.Duration {
	if r <= 0 {
		return d
	}
	neg := d < 0
	if neg {
		d = -d
	}
	if m := d % r; m+m < r {
		d = d - m
	} else {
		d = d + r - m
	}
	if neg {
		return -d
	}
	return d
}

func toString(k kokizami.Kizamier) string {
	var stoppedAt string
	if k.StoppedAt().Unix() == 0 {
		stoppedAt = "*" + time.Now().In(time.Local).Format("2006-01-02 15:04:05")
	} else {
		stoppedAt = k.StoppedAt().In(time.Local).Format("2006-01-02 15:04:05")
	}

	return strconv.Itoa(k.ID()) + "\t" +
		k.Desc() + "\t" +
		k.StartedAt().In(time.Local).Format("2006-01-02 15:04:05") + "\t" +
		stoppedAt + "\t" +
		round(k.Elapsed(), time.Second).String()
}

// CmdStart starts a new task
// kokizami start [new desc]
func CmdStart(c *cli.Context) {
	args := c.Args()
	if len(args) == 0 {
		fp, err := ioutil.TempFile("", "tmp_")
		if err != nil {
			log.Println(err)
			return
		}
		defer func() {
			err := os.Remove(fp.Name())
			if err != nil {
				log.Println(err)
			}
		}()

		filepath := fp.Name()
		err = fp.Close()
		if err != nil {
			log.Println(err)
		}

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim"
		}
		cmd := exec.Command(editor, filepath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			log.Println(err)
			return
		}

		bytes, err := ioutil.ReadFile(filepath)
		if err != nil {
			log.Println(err)
			return
		}

		ss := strings.Split(string(bytes), string("\n"))
		if len(ss) < 1 {
			log.Println("invalid arguments. needs (desc, started_at, stopped_at)")
			return
		}
		k, err := kokizami.Start(ss[0])
		if err != nil {
			log.Println(err)
		}
		fmt.Println(toString(k))
	} else if len(args) == 1 {
		desc := args[0]
		k, err := kokizami.Start(desc)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(toString(k))
	} else {
		log.Println("start needs one arguments [desc]")
		return
	}
}

// CmdRestart starts a task from old task list
// kokizami restart [id]
func CmdRestart(c *cli.Context) {
	args := c.Args()
	if len(args) != 1 {
		log.Println("restart needs one arguments [id]")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		log.Println(err)
	}

	k, err := kokizami.Get(id)
	if err != nil {
		log.Println(err)
	}

	k, err = kokizami.Start(k.Desc())
	if err != nil {
		log.Println(err)
	}
	fmt.Println(toString(k))
}

// CmdEdit edits a specified task
// kokizami edit [id]
// kokizami edit [id] desc       [new desc]
// kokizami edit [id] started_at [new started_at]
// kokizami edit [id] stopped_at [new stopped_at]
func CmdEdit(c *cli.Context) {
	args := c.Args()
	if len(args) == 1 {
		// FIXME: very long method. make them shorten...
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Println(err)
			return
		}

		k, err := kokizami.Get(id)
		if err != nil {
			log.Println(err)
			return
		}

		fp, err := ioutil.TempFile("", "tmp_")
		if err != nil {
			log.Println(err)
			return
		}
		defer func() {
			err = os.Remove(fp.Name())
			if err != nil {
				log.Println(err)
				return
			}
		}()

		filepath := fp.Name()

		_, err = fp.WriteString(k.Desc() + "\n" +
			k.StartedAt().In(time.Local).Format("2006-01-02 15:04:05") + "\n" +
			k.StoppedAt().In(time.Local).Format("2006-01-02 15:04:05"))
		if err != nil {
			log.Println(err)
			return
		}
		err = fp.Close()
		if err != nil {
			log.Println(err)
			return
		}

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim"
		}
		cmd := exec.Command(editor, filepath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			log.Println(err)
			return
		}

		bytes, err := ioutil.ReadFile(filepath)
		if err != nil {
			log.Println(err)
			return
		}

		ss := strings.Split(string(bytes), string("\n"))
		if len(ss) < 3 {
			log.Println("invalid arguments. needs (desc, started_at, stopped_at)")
			return
		}
		// kokizami: fixme. should be done by one transaction
		k, err = kokizami.Edit(id, "desc", ss[0])
		if err != nil {
			log.Println(err)
			return
		}

		startedAt, err := time.ParseInLocation("2006-01-02 15:04:05", ss[1], time.Local)
		startedAtStr := startedAt.UTC().Format("2006-01-02 15:04:05")
		k, err = kokizami.Edit(id, "started_at", startedAtStr)
		if err != nil {
			log.Println(err)
			return
		}
		stoppedAt, err := time.ParseInLocation("2006-01-02 15:04:05", ss[2], time.Local)
		stoppedAtStr := stoppedAt.UTC().Format("2006-01-02 15:04:05")
		k, err = kokizami.Edit(id, "stopped_at", stoppedAtStr)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(toString(k))
		return
	} else if len(args) == 3 {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Println(err)
			return
		}
		field := args[1]
		newValue := args[2]

		k, err := kokizami.Edit(id, field, newValue)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(k)
	} else {
		log.Println("edit needs three arguments " +
			"(id, [desc|started_at|stopped_at], [new value])")
		return
	}
}

// CmdList shows kokizami list
// kokizami list
func CmdList(c *cli.Context) {
	l, err := kokizami.List()
	if err != nil {
		log.Println(err)
		return
	}

	if len(l) == 0 {
		fmt.Println("list is empty")
		return
	}

	for _, v := range l {
		fmt.Println(toString(v))
	}
}

// CmdStop update specified task's stopped_at
// kokizami stop      ... stop all tasks they don't have stopped_at
// kokizami stop [id] ... stop a task by specified id
func CmdStop(c *cli.Context) {
	args := c.Args()
	switch len(args) {
	case 0:
		err := kokizami.StopAll()
		if err != nil {
			log.Println(err)
		}
	case 1:
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Println(err)
			return
		}
		err = kokizami.Stop(id)
		if err != nil {
			log.Println(err)
		}
	default:
		log.Println("stop needs at most one arguments [id]")
		return
	}
}

// CmdDelete deletes specified task
// kokizami delete [id]
func CmdDelete(c *cli.Context) {
	args := c.Args()
	if len(args) != 1 {
		log.Println("delete needs one arguments [id]")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		log.Println(err)
		return
	}

	err = kokizami.Delete(id)
	if err != nil {
		log.Println(err)
	}
}
