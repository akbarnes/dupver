package cmd

import (
	"github.com/spf13/cobra"

	"github.com/akbarnes/dupver/src/dupver"
	"github.com/akbarnes/dupver/src/fancyprint"
)

// diffCmd represents the diff command
var prefsCmd = &cobra.Command{
	Use:   "prefs",
	Short: "Print global preferences",
	Long: `This will print the global preferences.

Global preferences currently includes the specified diff
tool and the default repository`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := dupver.Options{}
		fancyprint.Setup(Debug, Verbose, Quiet, Monochrome)
		// TODO: print the preferences
		// fmt.Println("Global preferences:")
		dupver.PrintCurrentPreferences(opts)
	},
}

func init() {
	rootCmd.AddCommand(prefsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// diffCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// diffCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
