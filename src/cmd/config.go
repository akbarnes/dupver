package cmd

import (
	"os"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/akbarnes/dupver/src/dupver"
	"github.com/akbarnes/dupver/src/fancyprint"
)

var ConfigRepo bool
var EditCfg bool

// diffCmd represents the diff command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print working directory configuration",
	Long: `This will print the configuration of the current project working directory.

Configuration includes the project name, current branch, associated repositories and
which repository is the default.`,
	Run: func(cmd *cobra.Command, args []string) {

		opts := dupver.Options{JsonOutput: JsonOutput}
		prefs, _ := dupver.ReadPrefs(opts)
		fancyprint.Setup(Debug, Verbose, Quiet, Monochrome)
		// TODO: print the preferences
		// fmt.Println("Global preferences:")
		cfg, err := dupver.ReadWorkDirConfig(WorkDirPath)

		if err != nil {
			// Todo: handle invalid configuration file
			fancyprint.Warn("Could not read configuration file. Has the project working directory been initialized?")
			os.Exit(1)
		}

		if EditCfg {
			cfgPath := filepath.Join(WorkDirPath, ".dupver", "config.toml")
			dupver.EditFile(cfgPath, prefs)
			return
		}

		if ConfigRepo {
			if JsonOutput {
				dupver.PrintWorkDirReposConfigAsJson(cfg, opts)
			} else {
				dupver.PrintWorkDirReposConfig(cfg, opts)
			}
			return
		}

		if len(args) == 1 {
			key := args[0]

			switch key {
			case "WorkDirName":
				dupver.MultiPrint(cfg.WorkDirName, opts)
			case "Branch":
				dupver.MultiPrint(cfg.Branch, opts)
			case "DefaultRepo":
				dupver.MultiPrint(cfg.DefaultRepo, opts)
			default:
				fancyprint.Warnf("Key %s doesn't exit in the working directory configuration.", key)

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
			case "WorkDirName":
				cfg.WorkDirName = val
			case "Branch":
				cfg.Branch = val
			case "DefaultRepo":
				cfg.DefaultRepo = val
			default:
				fancyprint.Warnf("Key %s doesn't exit in the working directory configuration.", key)

				if JsonOutput {
					dupver.PrintJson(nil)
				}

				return
			}

			dupver.SaveWorkDirConfig(WorkDirPath, cfg, true, opts)
		}

		if JsonOutput {
			dupver.PrintCurrentWorkDirConfigAsJson(WorkDirPath, opts)
		} else {
			dupver.PrintCurrentPreferences(opts)
			fmt.Println("")
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
	configCmd.Flags().BoolVarP(&ConfigRepo, "repos", "R", false, "configure repos associated with this working directory")
	configCmd.Flags().BoolVarP(&EditCfg, "edit", "e", false, "edit the working directory configuration in the specified editor")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// diffCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
