package cmd

import (
	// "fmt"

	"github.com/akbarnes/dupver/src/dupver"
	"github.com/spf13/cobra"
)

// repoAddCmd represents the repoAdd command
var repoAddCmd = &cobra.Command{
	Use:   "add [name] [path]",
	Short: "Add a repository the current working directory",
	Long: `This command will add an additional repository
to the current project working directory. It takes
two optional positionaly arguments for the repo 
name and repo path. These take precedence over the
global command line flags`,
	Run: func(cmd *cobra.Command, args []string) {
		repoName := RepoName
		repoPath := RepoPath

		if len(args) >= 1 {
			repoName = args[0]
		}

		if len(args) >= 2 {
			repoPath = args[1]
		}

		// TODO: Read repoPath from environment variable if empty
		opts := dupver.SetVerbosity(dupver.Options{Color: true}, Verbose, Quiet)

		if Monochrome || Quiet {
			opts.Color = false
		}

		dupver.AddRepoToWorkDir(WorkDirPath, repoName, repoPath, opts)
	},
}

func init() {
	repoCmd.AddCommand(repoAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// repoAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// repoAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
