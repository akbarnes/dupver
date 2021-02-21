	package cmd

import (
	"fmt"
	"os"
	// "log"
	"path/filepath"

	"github.com/akbarnes/dupver/src/dupver"
	"github.com/spf13/cobra"
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
		opts := dupver.SetVerbosity(dupver.Options{Color: true}, Debug, Verbose, Quiet)
		tagName := ""

		if err != nil {
			// Todo: handle invalid configuration file
			fmt.Println("Could not read configuration file. Has the project working directory been initialized?")
			os.Exit(0)
		}

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

		opts.WorkDirName = cfg.WorkDirName
		opts.RepoName = RepoName
		opts.RepoPath = RepoPath

		if opts.Verbosity >= 2 {
			fmt.Printf("opts:\n%+v\n\n", opts)
		}

		headPath := filepath.Join(opts.RepoPath, "branches", opts.WorkDirName, "main.toml")

		if opts.Verbosity >= 2 {
			fmt.Printf("Head path: %s\n\n", headPath)
		}

		var snapshotId string

		if len(args) >= 1 {
			snapshotId = dupver.GetFullSnapshotId(args[0], opts)
		} else {
			mySnapshot, err := dupver.LastSnapshot(opts)

			if err != nil {
				fmt.Println("No snapshots found in project working directory. Have you initialized and commited yet?")
				os.Exit(1)
			}
			snapshotId = mySnapshot.ID
		}

		if opts.Verbosity >= 2 {
			fmt.Printf("Commit ID: %s\n\n", snapshotId)
		}

		if Monochrome || Quiet {
			opts.Color = false
		}

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
