/*
Copyright © 2020 Art Barnes <art@pin3.io>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := dupver.SetVerbosity(dupver.Options{Color: true}, Verbose, Quiet)
		
		if Monochrome || Quiet {
			opts.Color = false
		}

		opts.RepoName = RepoName
		opts.RepoPath = RepoPath
		
		if len(args) >= 1 {
			commitFile := args[0]
			tarFile := commitFile
			containingFolder := filepath.Dir(commitFile)

			if !strings.HasSuffix(commitFile, "tar") {
				fmt.Printf("%s -> %s, %s\n", commitFile, containingFolder, commitFile)
				tarFile = CreateTar(containingFolder, commitFile, opts)

				if len(Message) == 0 {
					Message = filepath.Base(commitFile)

					if opts.Verbosity >= 1 {
						fmt.Printf("Message not specified, setting to: %s\n", Message)
					}
				}
			}

			myHead := dupver.CommitFile(tarFile, parentIds, Message, opts)
			headPath := filepath.Join(containingFolder, ".dupver", "head.toml")
			dupver.WriteHead(headPath, myHead, opts)
		} else {
			dir, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			containingFolder := filepath.Dir(dir)
			workdirFolder := filepath.Base(dir)
			fmt.Printf("%s -> %s, %s\n", dir, containingFolder, workdirFolder)

			if len(Message) == 0 {
				Message = workdirFolder

				if opts.Verbosity >= 1 {
					fmt.Printf("Message not specified, setting to: %s\n", Message)
				}				
			}


			tarFile := CreateTar(containingFolder, workdirFolder, opts)
			myHead := dupver.CommitFile(tarFile, parentIds, Message, opts)	
			headPath := filepath.Join(".dupver", "head.toml")
			dupver.WriteHead(headPath, myHead, opts)
		}
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// copyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// copyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
