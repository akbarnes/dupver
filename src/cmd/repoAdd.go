package cmd

import (
	"github.com/spf13/cobra"

	"github.com/akbarnes/dupver/src/dupver"
	"github.com/akbarnes/dupver/src/fancyprint"
)

var SetDefaultRepo bool

// repoAddCmd represents the repoAdd command
var repoAddCmd = &cobra.Command{
	Use:   "add [repo_path] [repo_name]",
	Short: "Add a repository the current working directory",
	Long: `This adds an additional repository
to the current project working directory. 

The first optional positional argument
allows for the repository path to be specified. The second optional
positional argument allows for the repository name to be specified.
These take precedence over the global command-line flags. While
the positional arguments are considered optional, if they are
ommitted the path and name must be specified by the global
command-line flags.`,
	Run: func(cmd *cobra.Command, args []string) {
		repoName := RepoName
		repoPath := RepoPath

		if len(args) >= 1 {
			repoPath = args[0]
		}

		if len(args) >= 2 {
			repoName = args[1]
		}

		// TODO: Read repoPath from environment variable if empty
		opts := dupver.Options{JsonOutput: JsonOutput}
		fancyprint.Setup(Debug, Verbose, Quiet, Monochrome)
		dupver.AddRepoToWorkDir(WorkDirPath, repoName, repoPath, SetDefaultRepo, opts)
	},
}

func init() {
	repoCmd.AddCommand(repoAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// repoAddCmd.PersistentFlags().String("foo", "", "A help for foo")
	repoCmd.PersistentFlags().BoolVarP(&SetDefaultRepo, "default", "D", false, "set as default repo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// repoAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
