package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pankona/kokizami"
	"github.com/urfave/cli"
)

// GlobalFlags can be used globally
var GlobalFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "verbose",
		Usage: "specify to enable verbose mode",
	},
}

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
	{
		Name:   "summary",
		Usage:  "show summary of specified month",
		Action: CmdSummary,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "m, month",
				Value: thisMonth,
				Usage: "specify year and month to show summary",
			},
		},
	},
	{
		Name:   "tags",
		Usage:  "show list of tags",
		Action: CmdTags,
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "id",
				Usage: "show tags of specified task",
			},
		},
	},
}

// CommandNotFound is called when specified subcommand is not found
func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.",
		c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}

var thisMonth = func() string {
	return time.Now().Format("2006-01")
}()

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

func toString(k *kokizami.Kizami) string {
	var stoppedAt string
	if k.StoppedAt.Unix() == 0 {
		stoppedAt = "*" + time.Now().In(time.Local).Format("2006-01-02 15:04:05")
	} else {
		stoppedAt = k.StoppedAt.In(time.Local).Format("2006-01-02 15:04:05")
	}

	return strconv.Itoa(k.ID) + "\t" +
		k.Desc + "\t" +
		k.StartedAt.In(time.Local).Format("2006-01-02 15:04:05") + "\t" +
		stoppedAt + "\t" +
		round(k.Elapsed(), time.Second).String()
}

func kkzm(c *cli.Context) *kokizami.Kokizami {
	k := c.App.Metadata["kkzm"].(*kokizami.Kokizami)
	k.EnableVerboseQuery(c.GlobalBool("verbose"))
	return k
}

// CmdStart starts a new task
// kokizami start [new desc]
func CmdStart(c *cli.Context) error {
	args := c.Args()

	var desc string
	if len(args) == 0 {
		filepath, err := editTextWithEditor("")
		if err != nil {
			return err
		}

		bytes, err := ioutil.ReadFile(filepath) // #nosec
		if err != nil {
			return err
		}

		ss := strings.Split(string(bytes), string("\n"))
		if len(ss) < 1 {
			return fmt.Errorf("invalid arguments. needs (desc, started_at, stopped_at)")
		}
		desc = ss[0]
	} else if len(args) == 1 {
		desc = args[0]
	} else {
		return fmt.Errorf("start needs one arguments [desc]")
	}

	kkzm := kkzm(c)

	k, err := kkzm.Start(desc)
	if err != nil {
		return err
	}
	fmt.Println(toString(k))

	return tagging(kkzm, k.ID, desc)
}

func tagging(kkzm *kokizami.Kokizami, kizamiID int, desc string) error {
	// remove all tags from specified kizami first
	err := kkzm.Untagging(kizamiID)
	if err != nil {
		return err
	}

	tags := extractTagsFromString(desc)
	if len(tags) == 0 {
		return nil
	}

	err = kkzm.AddTags(tags)
	if err != nil {
		return err
	}

	tagIDs := make([]int, len(tags))
	ts, err := kkzm.TagsByTags(tags)
	if err != nil {
		return err
	}
	for i, v := range ts {
		tagIDs[i] = v.ID
	}

	return kkzm.Tagging(kizamiID, tagIDs)
}

// CmdRestart starts a task from old task list
// kokizami restart [id]
func CmdRestart(c *cli.Context) error {
	args := c.Args()
	if len(args) != 1 {
		return fmt.Errorf("restart needs one arguments [id]")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	kkzm := kkzm(c)

	k, err := kkzm.Get(id)
	if err != nil {
		return err
	}

	k, err = kkzm.Start(k.Desc)
	if err != nil {
		return err
	}
	fmt.Println(toString(k))

	return tagging(kkzm, k.ID, k.Desc)
}

// CmdEdit edits a specified task
func CmdEdit(c *cli.Context) error {
	args := c.Args()

	switch len(args) {
	// len(args) == 1 means that the whole of task will be edited with text editor
	// e.g) kkzm edit [id]
	case 1:
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		k, err := editTaskWithEditor(kkzm(c), id)
		if err != nil {
			return err
		}

		fmt.Println(toString(k))
		return nil
	}

	return fmt.Errorf("edit needs three arguments (id, [desc|started_at|stopped_at], [new value])")
}

// CmdList shows kokizami list
// kokizami list
func CmdList(c *cli.Context) error {
	l, err := kkzm(c).List()
	if err != nil {
		return err
	}

	if len(l) == 0 {
		fmt.Println("list is empty")
		return nil
	}

	buf := bytes.NewBuffer([]byte{})
	for _, v := range l {
		fmt.Fprintln(buf, toString(&v))
	}
	fmt.Printf("%s", buf)
	return nil
}

// CmdStop update specified task's stopped_at
// kokizami stop      ... stop all tasks they don't have stopped_at
// kokizami stop [id] ... stop a task by specified id
func CmdStop(c *cli.Context) error {
	args := c.Args()
	switch len(args) {
	case 0:
		err := kkzm(c).StopAll()
		if err != nil {
			return err
		}
	case 1:
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		err = kkzm(c).Stop(id)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("stop needs at most one arguments [id]")
	}
	return nil
}

// CmdDelete deletes specified task
// kokizami delete [id]
func CmdDelete(c *cli.Context) error {
	args := c.Args()
	if len(args) != 1 {
		return fmt.Errorf("delete needs one arguments [id]")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	return kkzm(c).Delete(id)
}

func editTaskWithEditor(kkzm *kokizami.Kokizami, id int) (*kokizami.Kizami, error) {
	k, err := kkzm.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to task by ID: %v", err)
	}

	filename, err := editTextWithEditor(fmt.Sprintf("%s\n%s\n%s",
		k.Desc,
		k.StartedAt.In(time.Local).Format("2006-01-02 15:04:05"),
		k.StoppedAt.In(time.Local).Format("2006-01-02 15:04:05")))
	if err != nil {
		return nil, fmt.Errorf("failed to edit text with editor: %v", err)
	}
	defer func() {
		e := os.Remove(filename)
		if e != nil {
			fmt.Printf("%v\n", e)
		}
	}()

	bytes, err := ioutil.ReadFile(filename) // #nosec
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	ss := strings.Split(string(bytes), string("\n"))
	if len(ss) < 3 {
		return nil, fmt.Errorf("invalid arguments. needs (desc, started_at, stopped_at)")
	}

	k, err = edit(kkzm, k, id, ss[0], ss[1], ss[2])
	if err != nil {
		return nil, fmt.Errorf("failed to edit a task: %v", err)
	}
	return k, nil
}

func edit(kkzm *kokizami.Kokizami, k *kokizami.Kizami, id int, desc, start, stop string) (*kokizami.Kizami, error) {
	startedAt, err := time.ParseInLocation("2006-01-02 15:04:05", start, time.Local)
	if err != nil {
		return nil, err
	}

	stoppedAt, err := time.ParseInLocation("2006-01-02 15:04:05", stop, time.Local)
	if err != nil {
		return nil, err
	}

	k.ID = id
	k.Desc = desc
	k.StartedAt = startedAt
	k.StoppedAt = stoppedAt

	ret, err := kkzm.Edit(k)
	if err != nil {
		return nil, err
	}

	return ret, tagging(kkzm, k.ID, desc)
}

func editTextWithEditor(prewrite string) (string, error) {
	f, err := ioutil.TempFile("", "tmp_")
	if err != nil {
		return "", fmt.Errorf("failed to open temporary file: %v", err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}()

	_, err = f.WriteString(prewrite)
	if err != nil {
		return "", fmt.Errorf("failed to write string on temporary file: %v", err)
	}

	err = runEditor(f.Name())
	if err != nil {
		return "", fmt.Errorf("failed to run editor: %v", err)
	}

	return f.Name(), nil
}

func runEditor(filename string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	cmd := exec.Command(editor, filename) // #nosec
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type descSummary struct {
	desc        string
	descElapsed time.Duration
}

type tagSummary struct {
	tagElapsed    time.Duration
	descSummaries []*descSummary
}

type tagSummaries map[string]*tagSummary

func (s tagSummaries) String() string {
	keys := make([]string, len(s))
	var index int
	for k := range s {
		keys[index] = k
		index++
	}
	sort.Strings(keys)

	buf := bytes.NewBuffer([]byte{})
	for _, v := range keys {
		tag := "-- No tag --"
		if v != "" {
			tag = v
		}
		fmt.Fprintf(buf, "%s\t%s\n", tag, s[v].tagElapsed)

		for _, d := range s[v].descSummaries {
			fmt.Fprintf(buf, "  %s\t%s\n", d.desc, d.descElapsed)
		}
	}
	return buf.String()
}

// CmdSummary shows summary of elapsed time of specified month
func CmdSummary(c *cli.Context) error {
	yyyymm := c.String("month")
	tags, err := kkzm(c).SummaryByTag(yyyymm)
	if err != nil {
		return err
	}

	summaries := tagSummaries{}
	for _, v := range tags {
		summaries[v.Tag] = &tagSummary{
			tagElapsed: v.Elapsed,
		}
	}

	descs, err := kkzm(c).SummaryByDesc(yyyymm)
	if err != nil {
		return err
	}

	for _, v := range descs {
		summaries[v.Tag].descSummaries = append(summaries[v.Tag].descSummaries, &descSummary{
			desc:        v.Desc,
			descElapsed: v.Elapsed,
		})
	}

	fmt.Printf("Summary of %s\n%s\n", yyyymm, summaries)
	return nil
}

// CmdTags shows list of tags
func CmdTags(c *cli.Context) error {
	var (
		ts  []kokizami.Tag
		err error
	)

	kkzm := kkzm(c)
	id := c.Int("id")

	if id == 0 {
		ts, err = kkzm.Tags()
		if err != nil {
			return err
		}
	} else {
		ts, err = kkzm.TagsByKizamiID(id)
		if err != nil {
			return err
		}
	}

	buf := bytes.NewBuffer([]byte{})
	for _, v := range ts {
		fmt.Fprintln(buf, v.Tag)
	}

	fmt.Printf("%s", buf)
	return nil
}

func extractTagsFromString(s string) []string {
	ss := strings.Split(s, " ")
	var tags []string
	for _, v := range ss {
		if strings.HasPrefix(v, "#") && len(v) >= 2 {
			tags = append(tags, v)
		}
	}
	return tags
}
