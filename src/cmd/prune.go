package cmd

import (
	"github.com/akbarnes/dupver/src/fancyprint"
	"github.com/spf13/cobra"
)

// pruneCmd represents the prune command
var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fancyprint.Setup(Debug, Verbose, Quiet, Monochrome)
		fancyprint.Notice("prune called")
	},
}

func init() {
	rootCmd.AddCommand(pruneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pruneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pruneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
