package cmd

import (
	"fmt"
	"os"
	// "log"
	// "path/filepath"
	
	"github.com/spf13/cobra"

	"github.com/akbarnes/dupver/src/dupver"
	"github.com/akbarnes/dupver/src/fancyprint"
)

// statusCmd represents the status command
var branchCmd = &cobra.Command{
	Use:   "branch [branch_name]",
	Short: "Switch project working directory to a branch",
	Long:  `This will switch a project working directory to the specified branch. `,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := dupver.ReadWorkDirConfig(WorkDirPath)
		fancyprint.Setup(Debug, Verbose, Quiet, Monochrome)

		if err != nil {
			// Todo: handle invalid configuration file
			fancyprint.Warn("Could not read configuration file. Has the project working directory been initialized?")
			os.Exit(1)
		}

		fancyprint.Debugf("Workdir configuration: %+v\n", cfg)
		fancyprint.Debugf("\nOld name: %s\n", cfg.WorkDirName)
		branch := Branch

		if len(args) >= 1 {
			branch = args[0]

			if len(branch) > 0 {
				cfg.Branch = branch
			}

			fancyprint.Debugf("\nNew name: %s\n", cfg.WorkDirName)

			cfg.Save(WorkDirPath, true)
		} else {
			if fancyprint.Verbosity >= fancyprint.NoticeLevel {
				fmt.Printf("Current branch: %s\n", cfg.Branch)
			} else {
				fmt.Println(cfg.Branch)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
