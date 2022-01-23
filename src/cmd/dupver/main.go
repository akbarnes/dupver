package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/akbarnes/dupver"
)

var OutputFolder string

func AddOptionFlags(fs *flag.FlagSet) {
	fs.BoolVar(&dupver.VerboseMode, "verbose", false, "verbose mode")
	fs.BoolVar(&dupver.VerboseMode, "v", false, "verbose mode")
	fs.BoolVar(&dupver.DebugMode, "debug", false, "debug mode")
	fs.BoolVar(&dupver.DebugMode, "d", false, "debug mode")
	fs.BoolVar(&dupver.QuietMode, "quiet", false, "quiet mode")
	fs.BoolVar(&dupver.QuietMode, "q", false, "quiet mode")
	fs.BoolVar(&dupver.RandomPoly, "random-poly", false, "generate random polynomial on initialization")
	fs.BoolVar(&dupver.RandomPoly, "r", false, "generate random polynomial on initialization")
}

func PostProcessOptionFlags() {
    if dupver.DebugMode {
        dupver.VerboseMode = true
    }
}

func main() {
    prefs, err := dupver.ReadPrefs(true)

    if err != nil {
        fmt.Fprintf(os.Stderr, "Error writing default prefs\n")
    }

	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	commitCmd := flag.NewFlagSet("commit", flag.ExitOnError)
	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)
	diffCmd := flag.NewFlagSet("diff", flag.ExitOnError)
	logCmd := flag.NewFlagSet("log", flag.ExitOnError)
	checkoutCmd := flag.NewFlagSet("checkout", flag.ExitOnError)
	repackCmd := flag.NewFlagSet("repack", flag.ExitOnError)

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Expected subcommand\n")
		os.Exit(1)
	}

	cmd := os.Args[1]

	if cmd == "init" {
		AddOptionFlags(initCmd)
		initCmd.Parse(os.Args[2:])
        PostProcessOptionFlags()
		dupver.ReadRepoConfig(true)
	} else if cmd == "commit" || cmd == "ci" {
		message := ""
		AddOptionFlags(commitCmd)
		commitCmd.Parse(os.Args[2:])
        PostProcessOptionFlags()
		cfg, err := dupver.ReadRepoConfig(true)
		filters, err := dupver.ReadFilters()

		if err != nil {
			// TODO: write to stderr
			fmt.Fprintf(os.Stderr, "Couldn't read filters file\n")
			os.Exit(1)
		}

		if commitCmd.NArg() >= 1 {
			message = commitCmd.Arg(0)
		}

		dupver.CommitSnapshot(message, filters, cfg.ChunkerPoly, cfg.PackSize, cfg.CompressionLevel)
	} else if cmd == "status" || cmd == "st" {
		dupver.AbortIfIncorrectRepoVersion()
		AddOptionFlags(statusCmd)
		statusCmd.Parse(os.Args[2:])
        PostProcessOptionFlags()
		filters, err := dupver.ReadFilters()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Encountered error when trying to read filters file, aborting:\n%v\n", err)
		}

		if statusCmd.NArg() >= 1 {
			dupver.DiffSnapshot(statusCmd.Arg(0), filters)
		} else {
			dupver.DiffSnapshot("", filters)
		}
	} else if cmd == "diff" {
		dupver.AbortIfIncorrectRepoVersion()
		AddOptionFlags(diffCmd)
		diffCmd.Parse(os.Args[2:])
        PostProcessOptionFlags()

		if diffCmd.NArg() >= 1 {
			dupver.DiffToolSnapshotFile(diffCmd.Arg(0), prefs.DiffTool)
		} else {
			dupver.DiffToolSnapshot(prefs.DiffTool)
        }
	} else if cmd == "log" {
		dupver.AbortIfIncorrectRepoVersion()
		AddOptionFlags(logCmd)
		logCmd.Parse(os.Args[2:])
        PostProcessOptionFlags()

		if logCmd.NArg() >= 1 {
			dupver.LogSingleSnapshot(logCmd.Arg(0))
		} else {
			dupver.LogAllSnapshots()
		}
	} else if cmd == "checkout" || cmd == "co" {
		dupver.AbortIfIncorrectRepoVersion()
		AddOptionFlags(checkoutCmd)
		checkoutCmd.StringVar(&OutputFolder, "out", "", "output folder")
		checkoutCmd.StringVar(&OutputFolder, "o", "", "output folder")
		checkoutCmd.Parse(os.Args[2:])
        PostProcessOptionFlags()

        checkoutFilter := "*"

        if checkoutCmd.NArg() >= 2 {
            checkoutFilter = checkoutCmd.Arg(1)
        }

		dupver.CheckoutSnapshot(checkoutCmd.Arg(0), OutputFolder, checkoutFilter)
	} else if cmd == "repack" {
		cfg, err := dupver.ReadRepoConfig(false)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't read repo configuration, exiting")
			os.Exit(1)
		}

		dupver.AbortIfIncorrectRepoVersion()
		AddOptionFlags(repackCmd)
		repackCmd.Parse(os.Args[2:])
        PostProcessOptionFlags()
        dupver.Repack(cfg.PackSize, cfg.CompressionLevel)
	} else if cmd == "version" || cmd == "ver" {
		fmt.Printf("%d.%d.%d\n", dupver.MajorVersion, dupver.MinorVersion, dupver.PatchVersion)
	} else {
        fmt.Fprintf(os.Stderr, "Unknown subcommand\n")
		os.Exit(1)
	}
}
