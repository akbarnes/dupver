package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// snapshots:
//	 project1:
//     snap1.json
//     snap2.json	 
// branches:
//   project1:
//     main.toml
//     feature.toml
// trees:
//   tree1.json
//   tree2.json
// packs:
//   fe:
//     fe1.zip
//     fe2.zip

// To copy:
// 1. copy snapshots folder directly
// 2. copy branch to reponame.branch
// 3. for packs:
//    a. don't copy duplicate chunks
//    b. copy new chunks into new packs
// 4. for trees, create new tree	
//    as new packs are created

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := dupver.SetVerbosity(dupver.Options{Color: true}, Verbose, Quiet)
		
		if Monochrome || Quiet {
			opts.Color = false
		}

		opts.RepoName = RepoName
		opts.RepoPath = RepoPath

		opts.DestRepoName := ""
		opts.DestRepoPath := destRepoPath

		
		if len(args) >= 1 {
			opts.DestRepoName := args[0]
		} 

		snapshotId = ""

		if len(args) >= 2 {
			snapshotId = args[2]
		}

		dupver.CopyRepo(snapshotId, opts)
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// copyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// copyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
