package cmd

import (
	// "fmt"
	"os"
	// "path/filepath"

	"github.com/spf13/cobra"

	"github.com/akbarnes/dupver/src/dupver"
	"github.com/akbarnes/dupver/src/fancyprint"
)

var SnapshotFiles bool

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log [commit_id]",
	Short: "List commits for the current working directory",
	Long: `This prints a list of commits for the current working directory."

If an optional positional argument is provided, this will specify
a commit ID to print in additional detail.`,
	Run: func(cmd *cobra.Command, args []string) {
		fancyprint.Setup(Debug, Verbose, Quiet, Monochrome)
		cfg, err := dupver.ReadWorkDirConfig(WorkDirPath)

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

		// Don't use LoadWorkDir so we don't load repo configs twice
		// if the repo name or path was changed via command line
		workDir := dupver.InstantiateWorkDir(cfg)

		if AllBranches {
			workDir.Branch = ""
		}

		// fancyprint.Debugf("Options: %+v\n", opts)
		fancyprint.Debugf("Workdir:%+v\n", workDir)
		snapshotId := ""
		fancyprint.Debugf("Commit ID: %s\n", snapshotId)

		// TODO: Yeesh...move this mess into a function
		if len(args) >= 1 {
			snapshotId = workDir.GetFullSnapshotId(args[0])
		
			fancyprint.Debugf("Full snapshot ID: %s\n\n", snapshotId)

			// Todo: fix specifying snapshot ID
			if JsonOutput {
				workDir.PrintSnapshotFilesAsJson(snapshotId)
			} else {
				workDir.PrintSnapshot(snapshotId)
			}	 
		} else {
			if JsonOutput {
				if SnapshotFiles {
					workDir.PrintSnapshotsAndFilesAsJson()
				} else {
					workDir.PrintSnapshotsAsJson()
				}
			} else {
				workDir.PrintSnapshots()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")
	rootCmd.PersistentFlags().BoolVarP(&SnapshotFiles, "files", "f", false, "Include files in JSON output")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
