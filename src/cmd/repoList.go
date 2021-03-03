package cmd

import (
	"github.com/spf13/cobra"

	"github.com/akbarnes/dupver/src/dupver"
	"github.com/akbarnes/dupver/src/fancyprint"
)

// repoListCmd represents the repoList command
var repoListCmd = &cobra.Command{
	Use:   "list",
	Short: "List repositories in the current working directory",
	Long: `This prints the repositories associated with the
current project working directory.

If the --quiet option is specified, then formatting will not
be applied and each row of the output will be the repo name
and repo path separated by a space.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Read repoPath from environment variable if empty
		opts := dupver.Options{}
		fancyprint.Setup(Debug, Verbose, Quiet, Monochrome)

		if JsonOutput {
			dupver.ListWorkDirReposAsJson(WorkDirPath, opts)
		} else {
			dupver.ListWorkDirRepos(WorkDirPath, opts)
		}
	},
}

func init() {
	repoCmd.AddCommand(repoListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// repoListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// repoListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
