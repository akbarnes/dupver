package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/akbarnes/dupver/src/dupver"	
	"github.com/akbarnes/dupver/src/fancyprint"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := dupver.ReadWorkDirConfig(WorkDirPath)
		opts := dupver.Options{}
		fancyprint.Setup(Debug, Verbose, Quiet, Monochrome)
		prefs, err := dupver.ReadPrefs()

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

		var snap dupver.Commit

		// snapshotId := dupver.GetFullSnapshotId(args[0], opts)
		if len(args) >= 1 {
			snapshotId := dupver.GetFullSnapshotId(args[0], opts)
			snap = dupver.ReadSnapshot(snapshotId, opts)
		} else {
			// TODO: check if err is not nil
			snap, err = dupver.LastSnapshot(opts)
		}

		randStr := dupver.RandHexString(40)
		tarFolder := filepath.Join(dupver.GetHome(), "temp", randStr)
		dupver.CreateFolder(tarFolder)
		tarPath := filepath.Join(tarFolder, "snapshot.tar")

		dupver.UnpackFile(tarPath, opts.RepoPath, snap.ChunkIDs)
		fancyprint.Debugf("Wrote to %s\n", tarPath)

		// TODO: Create a temporary folder to extract the tar file to
		tarCmd := exec.Command("tar", "xfv", tarPath)
		tarCmd.Dir = tarFolder
		fancyprint.Debugf("Running tar xfv %s", tarPath)	
		output, err := tarCmd.CombinedOutput()
	
		if err != nil {
			log.Fatal(fmt.Sprintf("Tar command failed\nOutput:\n%s\nError:\n%s\n", output, err))
		} else {
			fancyprint.Debugf("Ran tar command with output:\n%s\n", output)
		}

		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		containingFolder := filepath.Dir(dir)
		workdirFolder := filepath.Base(dir)

		committedFolder := filepath.Join(tarFolder, workdirFolder)

		diffCmd := exec.Command(prefs.DiffTool, committedFolder, workdirFolder)
		diffCmd.Dir = containingFolder
		diffCmd.Start()	
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// diffCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// diffCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
