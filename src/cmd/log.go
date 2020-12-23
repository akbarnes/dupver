package cmd

import (
	// "fmt"
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
		cfg = dupver.UpdateRepoPath(cfg, RepoPath)

		opts := dupver.SetVerbosity(dupver.Options{Color: true}, Verbose, Quiet)

		opts.RepoName = RepoName
		opts.RepoPath = RepoPath

		headPath := filepath.Join(WorkDirPath, ".dupver", "head.toml")
		myHead := dupver.ReadHead(headPath)
		snapshotId := myHead.CommitIDs[opts.RepoName]
		numSnapshots := 0

		// TODO: Yeesh...move this mess into a function
		if len(args) >= 1 {
			snapshotId = dupver.GetFullSnapshotId(args[0], cfg)
			numSnapshots = 1
		}

		
		if Monochrome || Quiet {
			opts.Color = false
		}
		
		dupver.PrintSnapshots(cfg, snapshotId, numSnapshots, opts)
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
