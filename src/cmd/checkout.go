package cmd

import (
	// "fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/akbarnes/dupver/src/dupver"	
	"github.com/akbarnes/dupver/src/fancyprint"
)

var OutFile string

// checkoutCmd represents the checkout command
var checkoutCmd = &cobra.Command{
	Use:   "checkout [commit_id]",
	Short: "Checkout commit to a tar file",
	Long: `This is used to restore a commit.

To avoid overwriting existing files (and because
the current architecture stores snapshots as a tar
file), the checkout command will export a commit to 
a tar file with the default name 
workdir_name-YYYY-MM-DDThh-mm-ss-commit_id[0:15].tar.
To specify a tar file name, use the --output flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := dupver.ReadWorkDirConfig(WorkDirPath)
		fancyprint.Setup(Debug, Verbose, Quiet, Monochrome)

		if err != nil {
			// Todo: handle invalid configuration file
			fancyprint.Warn("Could not read configuration file. Has the project working directory been initialized?")
			os.Exit(1)
		}

		if len(RepoName) > 0 {
			cfg.DefaultRepo = RepoName
		}

		// Don't use LoadWorkDir so we don't load repo configs twice
		// if the repo name or path was changed via command line
		workDir := dupver.InstantiateWorkDir(cfg)

		if len(RepoPath) > 0 {
			workDir.Repo.Path = RepoPath
			fancyprint.Debugf("Updating repo path to %s\n", RepoPath)
		}

		workDir.UnpackSnapshot(args[0], OutFile) 
	},
}

func init() {
	rootCmd.AddCommand(checkoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	checkoutCmd.Flags().StringVarP(&OutFile, "output", "o", "", "Output tar file")
}
