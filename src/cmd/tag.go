package cmd

import (
	"fmt"
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
		cfg := dupver.ReadWorkDirConfig(WorkDirPath)
		opts := dupver.SetVerbosity(dupver.Options{Color: true}, Debug, Verbose, Quiet)
		tagName := ""

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
			fmt.Println("opts:")
			fmt.Println(opts)
			fmt.Println("")
		}

		headPath := filepath.Join(opts.RepoPath, "branches", opts.WorkDirName, opts.BranchName+".toml")
		if opts.Verbosity >= 2 {
			fmt.Println("Head path:")
			fmt.Println(headPath)
			fmt.Println("")
		}

		myHead := dupver.ReadHead(headPath, opts)
		snapshotId := myHead.CommitID

		if opts.Verbosity >= 2 {
			fmt.Println("Commit ID:")
			fmt.Println(snapshotId)
			fmt.Println("")
		}

		if len(args) >= 1 {
			tagName = args[1]
		}

		if len(args) >= 2 {
			snapshotId = dupver.GetFullSnapshotId(args[2], opts)
		}

		if Monochrome || Quiet {
			opts.Color = false
		}

		dupver.CreateTag(tagName, snapshotId, opts)
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
