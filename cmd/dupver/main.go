package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/akbarnes/dupver"
	"github.com/restic/chunker"
)

var OutputFolder string

func AddOptionFlags(fs *flag.FlagSet) {
	fs.BoolVar(&dupver.VerboseMode, "verbose", false, "verbose mode")
	fs.BoolVar(&dupver.VerboseMode, "v", false, "verbose mode")
}

func main() {
	commitCmd := flag.NewFlagSet("commit", flag.ExitOnError)
	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)
	logCmd := flag.NewFlagSet("log", flag.ExitOnError)
	checkoutCmd := flag.NewFlagSet("checkout", flag.ExitOnError)

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Expected subcommand")
		os.Exit(1)
	}

	cmd := os.Args[1]

	if cmd == "init" {
		dupver.ReadRepoConfig(true)
	} else if cmd == "commit" || cmd == "ci" {
		cfg, err := dupver.ReadRepoConfig(true)

		message := ""
		AddOptionFlags(commitCmd)
		commitCmd.Parse(os.Args[2:])
		filters, err := dupver.ReadFilters()

		if err != nil {
			// TODO: write to stderr
			fmt.Println("Couldn't read filters file")
			os.Exit(1)
		}

		if commitCmd.NArg() >= 1 {
			message = commitCmd.Arg(0)
		}

		// TODO: need to use fixed poly
		rand.Seed(1)
		// p, _ := chunker.RandomPolynomial()
		const p chunker.Pol = 0x3abc9bff07d9e5

		if dupver.VerboseMode {
			fmt.Printf("Random polynomial: ")
			fmt.Println(p)
		}

		const packSize int64 = 500 * 1024 * 1024
		dupver.CommitSnapshot(message, filters, p, packSize, cfg.CompressionLevel)
	} else if cmd == "status" || cmd == "st" {
		dupver.AbortIfIncorrectRepoVersion()
		cfg, err := dupver.ReadRepoConfig(false)

		if err != nil {
			fmt.Printf("Can't read repo configuration, exiting")
			os.Exit(1)
		}

		AddOptionFlags(statusCmd)
		statusCmd.Parse(os.Args[2:])
		filters, err := dupver.ReadFilters()

		if err != nil {
			fmt.Printf("Encountered error when trying to read filters file, aborting:\n%v\n", err)
		}

		if statusCmd.NArg() >= 1 {
			dupver.DiffSnapshot(statusCmd.Arg(0), filters)
		} else {
			dupver.DiffSnapshot("", filters)
		}
	} else if cmd == "log" {
		dupver.AbortIfIncorrectRepoVersion()
		AddOptionFlags(logCmd)
		logCmd.Parse(os.Args[2:])

		if logCmd.NArg() >= 1 {
			snapshotNum, _ := strconv.Atoi(logCmd.Arg(0))
			dupver.LogSingleSnapshot(snapshotNum)
		} else {
			dupver.LogAllSnapshots()
		}
	} else if cmd == "checkout" || cmd == "co" {
		dupver.AbortIfIncorrectRepoVersion()
		AddOptionFlags(checkoutCmd)
		checkoutCmd.StringVar(&OutputFolder, "out", "", "output folder")
		checkoutCmd.StringVar(&OutputFolder, "o", "", "output folder")
		checkoutCmd.Parse(os.Args[2:])
		snapshotNum, _ := strconv.Atoi(checkoutCmd.Arg(0))
		dupver.CheckoutSnapshot(snapshotNum, OutputFolder)
	} else if cmd == "version" || cmd == "ver" {
		fmt.Printf("%d.%d.%f", dupver.DupverMajorversion, dupver.MinorVersion, dupver.PatchVersion)
	} else {
		fmt.Println("Unknown subcommand")
		os.Exit(1)
	}
}
