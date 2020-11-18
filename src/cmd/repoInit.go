package cmd

import (
	"github.com/akbarnes/dupver/src/dupver"
	"github.com/spf13/cobra"
)

var ChunkerPolynomial string

// initCmd represents the init command
var repoInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a repository",
	Long: `This initializes a dupver repository (as oppposed to a workdir)

This command can be run from anywhere but requires the repository path
as a positional argument. It takes an optional second positional
argument to specify the repository name.
Usage: dupver repo init <repopath> [<reponame>]`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := dupver.SetVerbosity(dupver.Options{Color: true}, Verbose, Quiet)

		if Monochrome || Quiet {
			opts.Color = false
		}

		repoName := RepoName
		repoPath := RepoPath

		if len(args) >= 1 {
			repoPath = args[0]
		}

		if len(args) >= 2 {
			repoName = args[1]
		}

		dupver.InitRepo(repoPath, repoName, ChunkerPolynomial, opts)
	},
}

func init() {
	repoCmd.AddCommand(repoInitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")
	repoInitCmd.PersistentFlags().StringVarP(&ChunkerPolynomial, "poly", "P", "", "specify chunker polynomial")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
