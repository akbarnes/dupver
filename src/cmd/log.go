package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/akbarnes/dupver/src/fancy_print"
	"github.com/akbarnes/dupver/src/dupver"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log [commit_id]",
	Short: "List commits for the current working directory",
	Long: `This prints a list of commits for the current working directory."

If an optional positional argument is provided, this will specify
a commit ID to print in additional detail.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := dupver.ReadWorkDirConfig(WorkDirPath)
		opts := dupver.SetVerbosity(dupver.Options{Color: true}, Debug, Verbose, Quiet)
		fancy_print.InitPrinting()
		fancy_print.SetVerbosityLevel(Debug, Verbose, Quiet)
		fancy_print.SetColoredOutput(Monochrome)

		if err != nil {
			// Todo: handle invalid configuration file
			fmt.Println("Could not read configuration file. Has the project working directory been initialized?")
			os.Exit(1)
		}


		if len(WorkDirPath) > 0 {
			os.Chdir(WorkDirPath)
		}

		if fancy_print.Verbosity >= fancy_print.InfoLevel {
			fmt.Println("cfg:")
			fmt.Println(cfg)
			fmt.Printf("\nRepo name: %s\nRepo path: %s\n\n", RepoName, RepoPath)
		}

		if len(RepoName) == 0 {
			RepoName = cfg.DefaultRepo
		}

		if len(RepoPath) == 0 {
			RepoPath = cfg.Repos[RepoName]

			if fancy_print.Verbosity >= fancy_print.InfoLevel {
				fmt.Printf("Updating repo path to %s\n", RepoPath)
			}
		}

		if len(Branch) == 0 {
			Branch = cfg.Branch
		}

		opts.WorkDirName = cfg.WorkDirName
		opts.RepoName = RepoName
		opts.RepoPath = RepoPath
		opts.Branch = Branch

		if AllBranches {
			opts.Branch = ""
		}

		if fancy_print.Verbosity >= fancy_print.InfoLevel {
			fmt.Println("opts:")
			fmt.Println(opts)
			fmt.Println("")
		}

		headPath := filepath.Join(opts.RepoPath, "branches", cfg.WorkDirName, "main.toml")

		if fancy_print.Verbosity >= fancy_print.InfoLevel {
			fmt.Println("Head path:")
			fmt.Println(headPath)
			fmt.Println("")
		}

		snapshotId := ""

		if fancy_print.Verbosity >= fancy_print.InfoLevel {
			fmt.Println("Commit ID:")
			fmt.Println(snapshotId)
			fmt.Println("")
		}

		// TODO: Yeesh...move this mess into a function
		if len(args) >= 1 {
			snapshotId = dupver.GetFullSnapshotId(args[0], opts)
		}

		if fancy_print.Verbosity >= fancy_print.InfoLevel {
			fmt.Printf("Full snapshot ID: %s\n", snapshotId)
		}

		if Monochrome || Quiet {
			opts.Color = false
		}

		dupver.PrintAllSnapshots(snapshotId, opts)
	},
}

func init() {
	rootCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
