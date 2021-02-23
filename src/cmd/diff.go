package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"log"
	"path/filepath"

	"github.com/akbarnes/dupver/src/dupver"
	"github.com/spf13/cobra"
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
		opts := dupver.SetVerbosity(dupver.Options{Color: true}, Debug, Verbose, Quiet)

		if err != nil {
			// Todo: handle invalid configuration file
			fmt.Println("Could not read configuration file. Has the project working directory been initialized?")
			os.Exit(1)
		}

		if Monochrome || Quiet {
			opts.Color = false
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

		var snap dupver.Commit

		// snapshotId := dupver.GetFullSnapshotId(args[0], opts)
		if len(args) >= 1 {
			snapshotId := dupver.GetFullSnapshotId(args[0], opts)
			snap = dupver.ReadSnapshot(snapshotId, opts)
		} else {
			// TODO: check if err is not nil
			snap, err = dupver.LastSnapshot(opts)
		}

		tarFile := dupver.RandHexString(40) + ".tar"
		tarFolder := filepath.Join(dupver.GetHome(), "temp")
		tarPath := filepath.Join(tarFolder, tarFile)

		dupver.UnpackFile(tarPath, opts.RepoPath, snap.ChunkIDs, opts)

		if opts.Verbosity >= 1 {
			fmt.Printf("Wrote to %s\n", tarPath)
		} else {
			fmt.Printf(tarPath)
		}

		// TODO: Create a temporary folder to extract the tar file to
		tarCmd := exec.Command("tar", "xfv", tarPath)
		tarCmd.Dir = tarFolder
	
		if opts.Verbosity >= 1 {
			log.Printf("Running tar xfv %s %s", tarPath)
		}
	
		output, err := tarCmd.CombinedOutput()
	
		if err != nil {
			log.Fatal(fmt.Sprintf("Tar command failed\nOutput:\n%s\nError:\n%s\n", output, err))
		} else if opts.Verbosity >= 3 {
			fmt.Printf("Ran tar command with output:\n%s\n", output)
		}

		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		containingFolder := filepath.Dir(dir)
		workdirFolder := filepath.Base(dir)

		committedFolder := filepath.Join(tarFolder, workdirFolder)

		diffCmd := exec.Command("bcomp.exe", committedFolder, workdirFolder)
		diffCmd.Dir = containingFolder
		diffCmd.Start()
	
		// if opts.Verbosity >= 1 {
		// 	log.Printf("Running tar cfv %s %s", tarPath, cleanCommitPath)
		// }
	
		// output, err := tarCmd.CombinedOutput()
	
		// if err != nil {
		// 	log.Fatal(fmt.Sprintf("Tar command failed\nOutput:\n%s\nError:\n%s\n", output, err))
		// } else if opts.Verbosity >= 3 {
		// 	fmt.Printf("Ran tar command with output:\n%s\n", output)
		// }	
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
