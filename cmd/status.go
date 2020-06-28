/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"

	"../dupver"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print file modification status of the project working directory",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("status called")
		snapshotId := ""

		if len(args) >= 1 {
			snapshotId = args[0]
		}

		cfg := dupver.ReadWorkDirConfig(WorkDirPath)
		cfg = dupver.UpdateRepoPath(cfg, RepoPath)
		var mySnapshot dupver.Commit
		
		if len(snapshotId) == 0 {
			snapshotPaths := dupver.ListSnapshots(cfg)
			mySnapshot = dupver.ReadSnapshotFile(snapshotPaths[len(snapshotPaths) - 1])
		} else {
			var err error
			mySnapshot, err = dupver.ReadSnapshotId(snapshotId, cfg)
			
			if err != nil {
				log.Fatal(fmt.Sprintf("Error reading snapshot %s", snapshotId))
			}
		}	

		verbosity := 1

		if Verbose {
			verbosity = 2
		} else if Quiet {
			verbosity = 0
		}

		dupver.WorkDirStatus(WorkDirPath, mySnapshot, verbosity)		
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
