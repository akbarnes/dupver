package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var UseDestPath bool

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
		cfg := dupver.ReadWorkDirConfig(WorkDirPath)
		cfg = dupver.UpdateRepoPath(cfg, RepoPath)

		sourcePath := cfg.RepoPath

		if len(RepoName) > 0 {
			sourcePath = cfg.Repos[RepoName]
		}

		if len(RepoPath) > 0 {
			sourcePath = RepoPath
		}

		destPath := args[0]

		if UseDestPath {
			destPath = cfg.repos[args[0]]
		} 		

		opts := dupver.SetVerbosity(dupver.Options{Color: true}, Verbose, Quiet)
		
		if Monochrome || Quiet {
			opts.Color = false
		}

		


		snapshotId = ""

		if len(args) >= 2 {
			snapshotId = args[2]
			dupver.CopySnapshot(snapshotId, sourcePath, destPath)
		} else {
			//dupver.CopyRepo(opts)
			fmt.Println("TODO: copy whole repo")
		}
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// copyCmd.PersistentFlags().String("foo", "", "A help for foo")
	rootCmd.PersistentFlags().BoolVarP(&UseDestPath, "path", "p", false, "specify destination repo path instead of name")		


	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// copyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
