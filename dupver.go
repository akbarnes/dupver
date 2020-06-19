package main


import (
	"flag"
	"fmt"
	// "os"
	// "github.com/google/subcommands"
)

const version string = "0.1.0-alpha"

func NewCommitCommand() *CommitCommand {
    cc := &CommitCommand{
        fs: flag.NewFlagSet("commit", flag.ContinueOnError),
    }

    cc.fs.StringVar(&cc.Message, "message", "", "Commit message")
    cc.fs.StringVar(&cc.Message, "m", "", "Commit message (shorthand)")

    return cc
}

type CommitCommand struct {
	fs *flag.FlagSet
	RepoPath string
    Message string
}

func (c *CommitCommand) Name() string {
    return c.fs.Name()
}

func (c *CommitCommand) Init(args []string) error {
    return c.fs.Parse(args)
}

func (c *CommitCommand) Run() error {
    fmt.Println("Hello", c.Message, "!")
    return nil
}

type InitCommand struct {
	fs *flag.FlagSet
	RepoPath string	
	ProjectDir string
	ProjectName string
}

func (c *InitCommand) Message() string {
    return c.fs.Message()
}

func (c *InitCommand) Init(args []string) error {
    return c.fs.Parse(args)
}

func (c *InitCommand) Run() error {
	if len(posArgs) >= 2 {
		workDir = posArgs[1]
	}

	// Read repoPath from environment variable if empty
	InitWorkDir(workDir, workDirName, repoPath)
	return nil
}

type Runner interface {
	Init([]string) error
	Run() error
    Message() string
}

func root(args []string) error {
    if len(args) < 1 {
        return errors.New("You must pass a sub-command")
    }

    cmds := []Runner{
        NewGreetCommand(),
    }

    subcommand := os.Args[1]

    for _, cmd := range cmds {
        if cmd.Name() == subcommand {
            cmd.Init(os.Args[2:])
            return cmd.Run()
        }
    }

    return fmt.Errorf("Unknown subcommand: %s", subcommand)
}

func main() {
    if err := root(os.Args[1:]); err != nil {
        fmt.Println(err)
        os.Exit(1)
	}
	
	// var filePath string
	// flag.StringVar(&filePath, "file", "", "Archive path")
	// flag.StringVar(&filePath, "f", "", "Archive path (shorthand)")

	// var repoPath string
	// flag.StringVar(&repoPath, "repository", "", "Repository path")
	// flag.StringVar(&repoPath, "r", "", "Repository path (shorthand)")

	// var workDir string
	// flag.StringVar(&workDir, "workdir", "", "Working directory")
	// flag.StringVar(&workDir, "d", "", "Working directory (shorthand)")

	// var workDirName string
	// flag.StringVar(&workDirName, "workdir-name", "", "Working directory name")
	// flag.StringVar(&workDirName, "w", "", "Working directory name (shorthand)")

	// var tagName string
	// flag.StringVar(&tagName, "tag-name", "", "Tag name")
	// flag.StringVar(&tagName, "t", "", "Tag name (shorthand)")

	// var verbosity int
	// flag.IntVar(&verbosity, "verbosity", 1, "Verbosity level")
	// flag.IntVar(&verbosity, "v", 1, "Verbosity level (shorthand)")	

	// flag.Parse()
	// posArgs := flag.Args()
	// cmd := posArgs[0]
	
  	// if cmd == "init-repo" {
	// 	repoPath := posArgs[1]
	// 	InitRepo(repoPath)
	// } else if cmd == "init" {
	// 	if len(posArgs) >= 2 {
	// 		workDir = posArgs[1]
	// 	}

	// 	// Read repoPath from environment variable if empty
	// 	InitWorkDir(workDir, workDirName, repoPath)
	// } else if cmd == "commit" || cmd == "ci" {
	// 	commitFile := posArgs[1]
	// 	CommitFile(commitFile, msg, verbosity)
	// } else if cmd == "checkout" || cmd == "co" {
	// 	snapshotId := posArgs[1]

	// 	cfg := ReadWorkDirConfig(workDir)
	// 	cfg = UpdateWorkDirName(cfg, workDirName)
	// 	cfg = UpdateRepoPath(cfg, repoPath)
	// 	snap := ReadSnapshot(snapshotId, cfg)

	// 	if len(filePath) == 0 {
	// 		timeStr := TimeToPath(snap.Time)
	// 		filePath = fmt.Sprintf("%s-%s-%s.tar", cfg.WorkDirName, timeStr, snapshotId[0:16])
	// 	}

	// 	UnpackFile(filePath, cfg.RepoPath, snap.ChunkIDs, verbosity) 
	// 	fmt.Printf("Wrote to %s\n", filePath)
	// } else if cmd == "log" || cmd == "list" {
	// 	snapshotId := ""

	// 	if len(posArgs) >= 2 {
	// 		snapshotId = posArgs[1]
	// 	}

	// 	cfg := ReadWorkDirConfig(workDir)
	// 	cfg = UpdateWorkDirName(cfg, workDirName)
	// 	cfg = UpdateRepoPath(cfg, repoPath)
	// 	PrintSnapshots(ListSnapshots(cfg), snapshotId)
	// } else if cmd == "version" {
	// 	fmt.Println("Dupver version:", version)
	// } else {
	// 	fmt.Println("No command specified, exiting")
	// 	fmt.Println("For available commands run: dupver -help")
	// }
}
