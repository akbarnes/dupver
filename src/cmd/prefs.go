package cmd

import (
	"os"
	
	"github.com/spf13/cobra"

	"github.com/akbarnes/dupver/src/dupver"
	"github.com/akbarnes/dupver/src/fancyprint"
)

var EditPrefs bool

// diffCmd represents the diff command
var prefsCmd = &cobra.Command{
	Use:   "prefs",
	Short: "Print global preferences",
	Long: `This will print the global preferences.

Global preferences currently includes the specified diff
tool and the default repository`,
	Run: func(cmd *cobra.Command, args []string) {
		fancyprint.Setup(Debug, Verbose, Quiet, Monochrome)
		prefs, err := dupver.ReadPrefs()
		// TODO: print the preferences
		// fmt.Println("Global preferences:")

		if err != nil {
			fancyprint.Warn("Could not read preferences file.")
			os.Exit(1)
		}		

		if EditPrefs {
			dupver.EditFile(dupver.GetPrefsPath(), prefs)
			return
		}

		if len(args) == 1 {
			key := args[0]

			if key == "Editor" || key == "editor" {
				dupver.MultiPrint(prefs.Editor, JsonOutput)
			} else if key == "DiffTool" || key == "difftool" {
				dupver.MultiPrint(prefs.DiffTool, JsonOutput)
			} else if key == "DefaultRepo" || key == "defaultrepo" {
				dupver.MultiPrint(prefs.DefaultRepo, JsonOutput)
			} else {
				fancyprint.Warnf("Key %s doesn't exit in the global preferences.", key)

				if JsonOutput {
					dupver.PrintJson(nil)
				}

				return
			}
		} else if len(args) >= 2 {
			key := args[0]
			val := args[1]

			if key == "Editor" || key == "editor" {
				prefs.Editor = val
			} else if key == "DiffTool" || key == "difftool" || key == "diff" {
				prefs.DiffTool = val
			} else if key == "DefaultRepo" || key == "defaultrepo" || key == "repo" {
				prefs.DefaultRepo = val
			} else {
				fancyprint.Warnf("Key %s doesn't exit in the global preferences.", key)

				if JsonOutput {
					dupver.PrintJson(nil)
				}

				return
			}

			prefs.Save(true)
		}

		if JsonOutput {
			prefs.PrintJson()
		} else {
			prefs.Print()
		}
	},
}

func init() {
	rootCmd.AddCommand(prefsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// diffCmd.PersistentFlags().String("foo", "", "A help for foo")
	prefsCmd.Flags().BoolVarP(&EditPrefs, "edit", "e", false, "edit the global preferences in the specified editor")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// diffCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
