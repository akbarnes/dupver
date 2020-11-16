package cmd

import (
	"github.com/akbarnes/dupver/src/dupver"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var repoInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a repository",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := dupver.SetVerbosity(dupver.Options{Color: true}, Verbose, Quiet)
		
		if Monochrome || Quiet {
			opts.Color = false
		}
				
		repoName := "main"
		repoPath := RepoPath

		if len(args) >= 1 {
			repoPath = args[1]
		}

		if len(args) >= 2 {
			repoName = args[1]
		}

		dupver.InitRepo(repoPath, repoName, opts)
	},
}

func init() {
	repoCmd.AddCommand(repoInitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
