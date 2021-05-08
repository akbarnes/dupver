package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/akbarnes/dupver/src/dupver"
	"github.com/akbarnes/dupver/src/fancyprint"
)

var ChunkerPolynomial string
var CopyRepoConfig string

// initCmd represents the init command
var repoInitCmd = &cobra.Command{
	Use:   "init path [repo_path] [repo_name]",
	Short: "Initialize a repository",
	Long: `This initializes a dupver repository (as oppposed to a workdir)

If no arguments are provided this command will create a repository in
the current working directory. The first optional positional argument
allows for the repository path to be specified. The second optional
positional argument allows for the repository name to be specified.
if no repository name is specified, the repository takes on the default
name of "main."`,
	Run: func(cmd *cobra.Command, args []string) {
		fancyprint.Setup(Debug, Verbose, Quiet, Monochrome)
		repoName := RepoName
		repoPath := RepoPath

		if len(args) >= 1 {
			repoPath = args[0]
		}

		if len(args) >= 2 {
			repoName = args[1]
		}

		dupver.InitRepo(repoPath, repoName, ChunkerPolynomial, CompressionLevel, JsonOutput)

		if len(CopyRepoConfig) > 0 {
			fancyprint.Infof("Copying repo configuration from %s", CopyRepoConfig)
			cfg, err := dupver.ReadRepoConfig(CopyRepoConfig)

			if err != nil {
				// Todo: handle invalid configuration file
				fancyprint.Warn("Could not read repository configuration file. Has the repository been initialized?")
				os.Exit(1)
			}

			cfg.Save(repoPath, true)
		}
	},
}

func init() {
	repoCmd.AddCommand(repoInitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")
	repoInitCmd.PersistentFlags().StringVarP(&ChunkerPolynomial, "poly", "P", "", "specify chunker polynomial")
	repoInitCmd.PersistentFlags().StringVarP(&CopyRepoConfig, "copy", "c", "", "specify repo path to copy configuration from")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
