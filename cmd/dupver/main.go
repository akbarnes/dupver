package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/akbarnes/dupver"
	"github.com/restic/chunker"
)

func ReadFilters() []string {
	filterPath := ".dupver_ignore"
	var filters []string
	f, _ := os.Open(filterPath)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		filters = append(filters, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return filters
}

var Message string
var OutputFolder string

func AddOptionFlags(fs *flag.FlagSet) {
	fs.BoolVar(&gover.VerboseMode, "verbose", false, "verbose mode")
	fs.BoolVar(&gover.VerboseMode, "v", false, "verbose mode")
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

	if cmd == "commit" || cmd == "ci" {
		AddOptionFlags(commitCmd)
		commitCmd.Parse(os.Args[2:])
		filters := ReadFilters()

		if commitCmd.NArg() >= 1 {
			Message = commitCmd.Arg(0)
		}

		// TODO: need to use fixed poly
		rand.Seed(1)
		// p, _ := chunker.RandomPolynomial()
		const p chunker.Pol = 0x3abc9bff07d9e5

		if gover.VerboseMode {
			fmt.Printf("Random polynomial: ")
			fmt.Println(p)
		}

		const packSize int64 = 100 * 1024 * 1024
		gover.CommitSnapshot(Message, filters, p, packSize)
	} else if cmd == "status" || cmd == "st" {
		AddOptionFlags(statusCmd)
		statusCmd.Parse(os.Args[2:])
		filters := ReadFilters()

		if statusCmd.NArg() >= 1 {
			gover.DiffSnapshot(statusCmd.Arg(0), filters)
		} else {
			gover.DiffSnapshot("", filters)
		}
	} else if cmd == "log" {
		AddOptionFlags(logCmd)
		logCmd.Parse(os.Args[2:])

		if logCmd.NArg() >= 1 {
			snapshotNum, _ := strconv.Atoi(logCmd.Arg(0))
			gover.LogSingleSnapshot(snapshotNum)
		} else {
			gover.LogAllSnapshots()
		}
	} else if cmd == "checkout" || cmd == "co" {
		AddOptionFlags(checkoutCmd)
		checkoutCmd.StringVar(&OutputFolder, "out", "", "output folder")
		checkoutCmd.StringVar(&OutputFolder, "o", "", "output folder")
		checkoutCmd.Parse(os.Args[2:])
		snapshotNum, _ := strconv.Atoi(checkoutCmd.Arg(0))
		gover.CheckoutSnaphot(snapshotNum, OutputFolder)
	} else {
		fmt.Println("Unknown subcommand")
		os.Exit(1)
	}
}
