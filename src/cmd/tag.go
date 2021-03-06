	package cmd

import (
	// "fmt"
	"os"
	// "log"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/akbarnes/dupver/src/dupver"	
	"github.com/akbarnes/dupver/src/fancyprint"
)

// statusCmd represents the status command
var tagCmd = &cobra.Command{
	Use:   "tag [tag_name] [commit_id]",
	Short: "Manage tags for commits",
	Long: `This will print tags, tag a commit, or delete tags.

Without any arguments this will list tags for the repository. If a
tag name and commit id are provided, this will add a tag for the 
specifed commit id. If only a tag is specified `,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := dupver.ReadWorkDirConfig(WorkDirPath)
		opts := dupver.Options{}
		tagName := ""

		if err != nil {
			// Todo: handle invalid configuration file
			fancyprint.Warn("Could not read configuration file. Has the project working directory been initialized?")
			os.Exit(0)
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

		opts.WorkDirName = cfg.WorkDirName
		opts.RepoName = RepoName
		opts.RepoPath = RepoPath

		fancyprint.Debugf("Options: %+v\n", opts)
		headPath := filepath.Join(opts.RepoPath, "branches", opts.WorkDirName, "main.toml")
		fancyprint.Debugf("Head path: %s\n", headPath)

		var snapshotId string

		if len(args) >= 1 {
			snapshotId = dupver.GetFullSnapshotId(args[0], opts)
		} else {
			mySnapshot, err := dupver.LastSnapshot(opts)

			if err != nil {
				fancyprint.Warn("No snapshots found in project working directory. Have you initialized and commited yet?")
				os.Exit(0)
			}
			snapshotId = mySnapshot.ID
		}

		fancyprint.Debugf("Commit ID: %s\n", snapshotId)

		if len(args) >= 1 {
			tagName = args[0]

			if len(args) >= 2 {
				snapshotId = dupver.GetFullSnapshotId(args[1], opts)
			}

			dupver.CreateTag(tagName, snapshotId, opts)
		}
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
