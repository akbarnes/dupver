package cmd

import (
	"fmt"

	"github.com/akbarnes/dupver/src/dupver"
	"github.com/spf13/cobra"
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
		opts := dupver.SetVerbosity(dupver.Options{Color: true}, Debug, Verbose, Quiet)

		if Monochrome || Quiet {
			opts.Color = false
		}

		cfg := dupver.ReadWorkDirConfig(WorkDirPath)
		cfg = dupver.UpdateRepoPath(cfg, RepoPath)
		snapshotId := dupver.GetFullSnapshotId(args[0], cfg)
		snap := dupver.ReadSnapshot(snapshotId, cfg)

		if len(OutFile) == 0 {
			timeStr := dupver.TimeToPath(snap.Time)
			OutFile = fmt.Sprintf("%s-%s-%s.tar", cfg.WorkDirName, timeStr, snap.ID[0:16])
		}

		dupver.UnpackFile(OutFile, cfg.RepoPath, snap.ChunkIDs, opts)

		if opts.Verbosity >= 1 {
			fmt.Printf("Wrote to %s\n", OutFile)
		} else {
			fmt.Printf(OutFile)
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
