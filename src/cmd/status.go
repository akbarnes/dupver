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
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := dupver.ReadWorkDirConfig(WorkDirPath)
		cfg = dupver.UpdateRepoPath(cfg, RepoPath)

		headPath := filepath.Join(WorkDirPath, ".dupver", "head.toml")
		myHead := dupver.ReadHead(headPath)
		snapshotId := myHead.CommitID

		if len(args) >= 1 {
			snapshotId = dupver.GetFullSnapshotId(args[0], cfg)
		}

		mySnapshot := dupver.ReadSnapshot(snapshotId, cfg)
		opts := dupver.SetVerbosity(dupver.Options{Color: true}, Verbose, Quiet)

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
