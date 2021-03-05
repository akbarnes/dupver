package cmd

import (
	// "fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/akbarnes/dupver/src/dupver"	
	"github.com/akbarnes/dupver/src/fancyprint"
)

var UseDestPath bool
var CopyAllProjects bool

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

		if len(Branch) == 0 {
			Branch = cfg.Branch
		}

		opts.WorkDirName = cfg.WorkDirName
		opts.RepoName = RepoName
		opts.RepoPath = RepoPath
		opts.Branch = Branch

		if AllBranches {
			opts.Branch = ""
		}

		sourcePath := opts.RepoPath

		if len(RepoPath) > 0 {
			sourcePath = RepoPath
		}

		destPath := args[0]

		if !UseDestPath {
			destPath = cfg.Repos[args[0]]
		}

		fancyprint.Infof("Source path: %s\n", sourcePath)
		fancyprint.Infof("Dest path: %s\n", destPath)		

		snapshotId := ""
		// TODO: look up the snapshot id based on 1st n characters

		if len(args) >= 2 {
			snapshotId = dupver.GetFullSnapshotId(args[1], opts)
			dupver.CopySnapshot(snapshotId, sourcePath, destPath, opts)
		} else {
			//dupver.CopyRepo(opts)
			snapshotPaths := dupver.ListSnapshots(opts)
			// print a specific revision
			for _, snapshotPath := range snapshotPaths {
				snap := dupver.ReadSnapshotFile(snapshotPath)

				// fmt.Printf("Path: %s\n", snapshotPath)
				fancyprint.Infof("Copying snapshot %s\n", snap.ID)
				branch := opts.Branch

				if len(branch) == 0 || len(branch) > 0 && branch == snap.Branch {
					dupver.CopySnapshot(snap.ID, sourcePath, destPath, opts)
				}
			}
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
	copyCmd.PersistentFlags().BoolVarP(&CopyAllProjects, "all-projects", "a", false, "copy all projects")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// copyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
