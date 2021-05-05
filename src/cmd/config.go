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
		fancyprint.Setup(Debug, Verbose, Quiet, Monochrome)
		prefs, _ := dupver.ReadPrefs()
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
				cfg.PrintReposAsJson()
			} else {
				cfg.PrintRepos()
			}
			return
		}

		if len(args) == 1 {
			key := args[0]

			if key == "WorkDirName" || key == "workdirname" || key == "workdir" || key == "name" {
				dupver.MultiPrint(cfg.WorkDirName, JsonOutput)
			} else if key == "Branch" || key == "branch" {
				dupver.MultiPrint(cfg.Branch, JsonOutput)
			} else if key == "DefaultRepo" || key == "defaultrepo" || key == "repo" {
				dupver.MultiPrint(cfg.DefaultRepo, JsonOutput)
			} else {
				fancyprint.Warnf("Key %s doesn't exit in the working directory configuration.", key)

				if JsonOutput {
					dupver.PrintJson(nil)
				}

				return
			}			
		} else if len(args) >= 2 {
			key := args[0]
			val := args[1]

			if key == "WorkDirName" || key == "workdirname" || key == "workdir" || key == "name" {
				cfg.WorkDirName = val
			} else if key == "Branch" || key == "branch" {
				cfg.Branch = val
			} else if key == "DefaultRepo" || key == "defaultrepo" || key == "repo" {
				cfg.DefaultRepo = val
			} else {
				fancyprint.Warnf("Key %s doesn't exit in the working directory configuration.", key)

				if JsonOutput {
					dupver.PrintJson(nil)
				}

				return
			}

			cfg.Save(WorkDirPath, true)
		}

		if JsonOutput {
			cfg.PrintJson()
		} else {
			prefs.Print()
			fmt.Println("")	
			cfg.Print()
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
