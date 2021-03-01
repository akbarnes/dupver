package cmd

import (
	// "fmt"
	"os"
	// "log"

	"github.com/spf13/cobra"

	"github.com/akbarnes/dupver/src/dupver"	
	"github.com/akbarnes/dupver/src/fancyprint"
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
		cfg, err := dupver.ReadWorkDirConfig(WorkDirPath)
		opts := dupver.Options{}
		fancyprint.Setup(Debug, Verbose, Quiet, Monochrome)

		if err != nil {
			// Todo: handle invalid configuration file
			fancyprint.Warn("Could not read configuration file. Has the project working directory been initialized?")
			os.Exit(1)
		}

		if len(WorkDirPath) > 0 {
			os.Chdir(WorkDirPath)
		}

		fancyprint.Debugf("Workdir Configuration: %+v\n", cfg)
		fancyprint.Debugf("Repo name: %s\nRepo path: %s\n", RepoName, RepoPath)

		if len(RepoName) == 0 {
			RepoName = cfg.DefaultRepo
		}

		if len(RepoPath) == 0 {
			RepoPath = cfg.Repos[RepoName]
			fancyprint.Debugf("Updating repo path to %s\n", RepoPath)
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

		fancyprint.Debugf("Options: %+v\n", opts)
		var mySnapshot dupver.Commit

		if len(args) >= 1 {
			snapshotId := dupver.GetFullSnapshotId(args[0], opts)
			mySnapshot = dupver.ReadSnapshot(snapshotId, opts)
		} else {
			mySnapshot, err = dupver.LastSnapshot(opts)

			if err != nil {
				fancyprint.Notice("No snapshots")
			}
		}

		fancyprint.Debugf("Snapshot commit ID: %s\n", mySnapshot.ID)
		dupver.WorkDirStatus("", mySnapshot, opts)
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
