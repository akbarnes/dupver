package cmd

import (
	// "fmt"

	"github.com/spf13/cobra"

	"github.com/akbarnes/dupver/src/dupver"	
	"github.com/akbarnes/dupver/src/fancyprint"
)

var ProjectName string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [project_working_dir]",
	Short: "Initialize a project working directory",
	Long: `This initializes a project working directory.

If an optional positional argument is provided, this will 
specify the location of the project working directory. 
Otherwise, the current working directory is used.`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := dupver.Options{}
		fancyprint.Setup(Debug, Verbose, Quiet, Monochrome)

		workDirPath := WorkDirPath

		if len(RepoName) == 0 {
			RepoName = "main"
		}

		if len(Branch) == 0 {
			Branch = "main"
		}

		opts.RepoName = RepoName
		opts.RepoPath = RepoPath
		opts.Branch = Branch

		if len(args) >= 1 {
			workDirPath = args[0]
		}

		// TODO: Read repoPath from environment variable if empty
		dupver.InitWorkDir(workDirPath, ProjectName, opts)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	initCmd.Flags().StringVarP(&ProjectName, "project-name", "p", "", "project name")
}
