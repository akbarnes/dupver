package cmd

import (
	"fmt"
	"path/filepath"

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
		cfg := dupver.ReadWorkDirConfig(WorkDirPath)
		opts := dupver.SetVerbosity(dupver.Options{Color: true}, Debug, Verbose, Quiet)
		// fmt.Println("Verbosity:")
		// fmt.Println(opts.Verbosity)
		// fmt.Println("")

		if opts.Verbosity >= 2 {
			fmt.Println("cfg:")
			fmt.Println(cfg)
			fmt.Printf("\nRepo name: %s\nRepo path: %s\n\n", RepoName, RepoPath)
		}

		if len(RepoName) == 0 {
			RepoName = cfg.DefaultRepo
		}

		if len(RepoPath) == 0 {
			RepoPath = cfg.Repos[RepoName]

			if opts.Verbosity >= 2 {
				fmt.Printf("Updating repo path to %s\n", RepoPath)
			}
		}

		if len(BranchName) == 0 {
			BranchName = cfg.BranchName
		}

		opts.WorkDirName = cfg.WorkDirName
		opts.RepoName = RepoName
		opts.RepoPath = RepoPath
		opts.BranchName = BranchName

		if opts.Verbosity >= 2 {
			fmt.Println("opts:")
			fmt.Println(opts)
			fmt.Println("")
		}

		headPath := filepath.Join(opts.RepoPath, "branches", cfg.WorkDirName, opts.BranchName + ".toml")

		if opts.Verbosity >= 2 {
			fmt.Println("Head path:")
			fmt.Println(headPath)
			fmt.Println("")
		}

		myHead := dupver.ReadHead(headPath)

		snapshotId := myHead.CommitID
		numSnapshots := 0

		if opts.Verbosity >= 2 {
			fmt.Println("Commit ID:")
			fmt.Println(snapshotId)
			fmt.Println("")
		}

		// TODO: Yeesh...move this mess into a function
		if len(args) >= 1 {
			snapshotId = dupver.GetFullSnapshotId(args[0], opts)
			numSnapshots = 1
		}

		
		if Monochrome || Quiet {
			opts.Color = false
		}
		
		dupver.PrintSnapshots(snapshotId, numSnapshots, opts)
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
