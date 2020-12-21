package cmd

import (
	"fmt"

	"github.com/akbarnes/dupver/src/dupver"
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
	Use:   "copy <dest> [snapshot_id]",
	Short: "Copy snapshots and chunks to the specified repo",
	Long: `This command will copy both snapshots and chunks from
either the default repo or a specified repo to the destination repo.

An alternate source repo can be specified with the --repo-name and
--repo-path flags. The destination repo is specified with the 
required first positional argument. By default, the destination
repo is specified from one of the repos current in the working 
directory configuration file. However, if it desired to specify
an arbitrary path that can be done by using the --path Boolean
flag, which instructs dupver to interpret the destination 
argument as a path. A second optional positional argument 
will limit only a single specified snapshot id to be copied.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := dupver.ReadWorkDirConfig(WorkDirPath)
		// cfg = dupver.UpdateRepoPath(cfg, RepoPath)

		sourcePath := cfg.RepoPath

		if len(RepoName) > 0 {
			cfg.RepoPath = cfg.Repos[RepoName]
		}

		if len(RepoPath) > 0 {
			sourcePath = RepoPath
		}

		destPath := args[0]

		if !UseDestPath {
			destPath = cfg.Repos[args[0]]
		} 		

		opts := dupver.SetVerbosity(dupver.Options{Color: true}, Verbose, Quiet)
		
		if Monochrome || Quiet {
			opts.Color = false
		}

		snapshotId := ""
		// TODO: look up the snapshot id based on 1st n characters

		if len(args) >= 2 {
			snapshotId = dupver.GetFullSnapshotId(args[1], cfg)
			dupver.CopySnapshot(cfg, snapshotId, sourcePath, destPath, opts)
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
	copyCmd.PersistentFlags().BoolVarP(&UseDestPath, "path", "p", false, "specify destination repo path instead of name")		


	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// copyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
