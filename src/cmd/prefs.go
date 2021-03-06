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
		prefs, _ := dupver.ReadPrefs(opts)
		// TODO: print the preferences
		// fmt.Println("Global preferences:")

		if len(args) == 1 {
			key := args[0]

			switch key {
			case "Editor":
				dupver.MultiPrint(prefs.Editor, opts)
			case "DiffTool":
				dupver.MultiPrint(prefs.DiffTool, opts)
			case "DefaultRepo":
				dupver.MultiPrint(prefs.DefaultRepo, opts)
			default:
				fancyprint.Warnf("Key %s doesn't exit in the global preferences.", key)

				if JsonOutput {
					dupver.PrintJson(nil)
				}
			}

			return
		}

		if len(args) >= 2 {
			key := args[0]
			val := args[1]

			switch key {
			case "Editor":
				prefs.Editor = val
			case "DiffTool":
				prefs.DiffTool = val
			case "DefaultRepo":
				prefs.DefaultRepo = val
			default:
				fancyprint.Warnf("Key %s doesn't exit in the global preferences.", key)

				if JsonOutput {
					dupver.PrintJson(nil)
				}

				return
			}

			dupver.SavePrefs(prefs, true, opts)
		}

		if JsonOutput {
			dupver.PrintCurrentPreferencesAsJson(opts)
		} else {
			dupver.PrintCurrentPreferences(opts)
		}
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
