package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/akbarnes/dupver/src/dupver"	
	"github.com/akbarnes/dupver/src/fancyprint"
)

var OutFile string

// checkoutCmd represents the checkout command
var checkoutCmd = &cobra.Command{
	Use:   "checkout [commit_id]",
	Short: "Checkout commit to a tar file",
	Long: `This is used to restore a commit.

To avoid overwriting existing files (and because
the current architecture stores snapshots as a tar
file), the checkout command will export a commit to 
a tar file with the default name 
workdir_name-YYYY-MM-DDThh-mm-ss-commit_id[0:15].tar.
To specify a tar file name, use the --output flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := dupver.ReadWorkDirConfig(WorkDirPath)
		opts := dupver.Options{}
		fancyprint.Setup(Debug, Verbose, Quiet, Monochrome)

		if err != nil {
			// Todo: handle invalid configuration file
			fancyprint.Warn("Could not read configuration file. Has the project working directory been initialized?")
			os.Exit(1)
		}

		if len(RepoName) == 0 {
			RepoName = cfg.DefaultRepo
		}

		if len(RepoPath) == 0 {
			RepoPath = cfg.Repos[RepoName]
			fancyprint.Debugf("Updating repo path to %s\n", RepoPath)
		}

		opts.WorkDirName = cfg.WorkDirName
		opts.RepoName = RepoName
		opts.RepoPath = RepoPath

		snapshotId := dupver.GetFullSnapshotId(args[0], opts)
		snap := dupver.ReadSnapshot(snapshotId, opts)

		if len(OutFile) == 0 {
			timeStr := dupver.TimeToPath(snap.Time)
			OutFile = fmt.Sprintf("%s-%s-%s.tar", opts.WorkDirName, timeStr, snap.ID[0:16])
		}

		dupver.UnpackFile(OutFile, opts.RepoPath, snap.ChunkIDs, opts)

		if fancyprint.Verbosity <= fancyprint.WarningLevel {
			fmt.Println(OutFile)
		} else {
			fmt.Printf("Wrote to %s\n", OutFile)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	checkoutCmd.Flags().StringVarP(&OutFile, "output", "o", "", "Output tar file")
}
