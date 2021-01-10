package cmd

import (
	"fmt"
	// "log"
	// "path/filepath"

	"github.com/akbarnes/dupver/src/dupver"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var branchCmd = &cobra.Command{
	Use:   "branch [branch_name]",
	Short: "Switch project working directory to a branch",
	Long: `This will switch a project working directory to the specified branch. `,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := dupver.ReadWorkDirConfig(WorkDirPath)
		opts := dupver.SetVerbosity(dupver.Options{Color: true}, Debug, Verbose, Quiet)

		if opts.Verbosity >= 2 {
			fmt.Println("cfg:")
			fmt.Println(cfg)
		}

		if opts.Verbosity >= 1 {
			fmt.Printf("\nOld name: %s\n", cfg.WorkDirName)
		}

		branch := Branch

		if len(args) >= 1 {
			branch = args[0]
		}

		if len(branch) > 0 {
			cfg.Branch = branch
		}

		if opts.Verbosity >= 1 {
			fmt.Printf("\nNew name: %s\n", cfg.WorkDirName)
		}		

		if Monochrome || Quiet {
			opts.Color = false
		}		

		dupver.SaveWorkDirConfig(WorkDirPath, cfg, true, opts)
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
