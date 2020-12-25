package cmd

import (
	// "fmt"
	// "log"
	"path/filepath"

	"github.com/akbarnes/dupver/src/dupver"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print file modification status of the project working directory",
	Long: `This will print the modification status of files in the current project working directory.

By default, new files are indicated with the prefix "+ " and are colored in green. Modified
files are indicated with the prefix "M " and are colored in cyan. Deleted files are
indicated with the prefix "- " and are colored in red. Dupver does not currently track file
renames (though this does not impact disk usage on account of the files being stored as
chunks rather than diffs). For usage as part of a command pipeline, colors can be disabled
with the --monochrome flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := dupver.ReadWorkDirConfig(WorkDirPath)
		cfg = dupver.UpdateRepoPath(cfg, RepoPath)

		opts := dupver.SetVerbosity(dupver.Options{Color: true}, Verbose, Quiet)
		opts.RepoName = RepoName
		opts.RepoPath = RepoPath

		headPath := filepath.Join(WorkDirPath, ".dupver", "head.toml")
		myHead := dupver.ReadHead(headPath)
		snapshotId := myHead.CommitIDs[opts.RepoName]

		if len(args) >= 1 {
			snapshotId = dupver.GetFullSnapshotId(args[0], cfg)
		}

		mySnapshot := dupver.ReadSnapshot(snapshotId, cfg)

		if Monochrome || Quiet {
			opts.Color = false
		}

		dupver.WorkDirStatus(WorkDirPath, mySnapshot, opts)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
