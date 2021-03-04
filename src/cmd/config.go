package cmd

import (
	"github.com/spf13/cobra"

	"github.com/akbarnes/dupver/src/dupver"
	"github.com/akbarnes/dupver/src/fancyprint"
)

// diffCmd represents the diff command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print working directory configuration",
	Long: `This will print the configuration of the current project working directory.

Configuration includes the project name, current branch, associated repositories and
which repository is the default.`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := dupver.Options{}
		fancyprint.Setup(Debug, Verbose, Quiet, Monochrome)
		// TODO: print the preferences
		// fmt.Println("Global preferences:")

		// fmt.Println("\nCurrent project working directory configuration:")
		if JsonOutput {
			dupver.PrintCurrentWorkDirConfigAsJson(WorkDirPath, opts)
		} else {
			dupver.PrintCurrentPreferences(opts)
			dupver.PrintCurrentWorkDirConfig(WorkDirPath, opts)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// diffCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// diffCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
