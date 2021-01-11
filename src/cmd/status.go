package cmd

import (
	"fmt"
	// "log"

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
chunks rather than diffs). For usage as part of a comm
	and pipeline, colors can be disabled
with the --monochrome flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := dupver.ReadWorkDirConfig(WorkDirPath)
		opts := dupver.SetVerbosity(dupver.Options{Color: true}, Debug, Verbose, Quiet)

		if opts.Verbosity >= 2 {
			fmt.Printf("cfg:\n%+v\n\n", cfg)
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

		opts.WorkDirName = cfg.WorkDirName
		opts.RepoName = RepoName
		opts.RepoPath = RepoPath
		opts.Branch = Branch		

		if AllBranches {
			opts.Branch = ""
		}

		if opts.Verbosity >= 2 {
			fmt.Printf("opts:\n%+v\n\n", opts)
		}

		var mySnapshot dupver.Commit

		if len(args) >= 1 {
			snapshotId := dupver.GetFullSnapshotId(args[0], opts)
			mySnapshot = dupver.ReadSnapshot(snapshotId, opts)
		} else {
			mySnapshot = dupver.LastSnapshot(opts)
		}

		if opts.Verbosity >= 2 {
			fmt.Printf("Snapshot commit ID: %s\n", mySnapshot.ID)
		}

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
