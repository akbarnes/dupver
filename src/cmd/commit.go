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
	"log"
	"path"
	"path/filepath"
	"strings"
	"os"
	"os/exec"

	"github.com/akbarnes/dupver/src/dupver"
	"github.com/spf13/cobra"
)

var Message string

func CreateTar(parentPath string, commitPath string, verbosity int) string {
	tarFile := dupver.RandHexString(40) + ".tar"
	tarFolder := path.Join(dupver.GetHome(), "temp")
	tarPath := path.Join(tarFolder, tarFile)

	// InitRepo(workDir)
	if verbosity >= 1 {
		fmt.Printf("Tar path: %s\n", tarPath)
		fmt.Printf("Creating folder %s\n", tarFolder)
	}

	os.Mkdir(tarFolder, 0777)	

	CompressTar(parentPath, commitPath, tarPath)
	return tarPath
}

func CompressTar(parentPath string, commitPath string, tarPath string) string {
	if len(tarPath) == 0 {
		tarPath = commitPath + ".tar"
	}

	cleanCommitPath := filepath.Clean(commitPath)

	tarCmd := exec.Command("tar", "cfv", tarPath, cleanCommitPath)
	tarCmd.Dir = parentPath
	log.Printf("Running tar cfv %s %s", tarPath, cleanCommitPath)
	output, err := tarCmd.CombinedOutput()	

	if err != nil {
		log.Fatal(fmt.Sprintf("Tar command failed\nOutput:\n%s\nError:\n%s\n", output, err))	
	} else {
		fmt.Printf("Ran tar command with output:\n%s\n", output)
	}

	return tarPath
}

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit a tar file into the repository",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		verbosity := dupver.SetVerbosity(Verbose, Quiet)

		if len(args) >= 1 {
			commitFile := args[0]
			tarFile := commitFile

			if !strings.HasSuffix(commitFile, "tar") {
				containingFolder := filepath.Dir(commitFile)
				fmt.Printf("%s -> %s, %s\n", commitFile, containingFolder, commitFile)
				tarFile = CreateTar(containingFolder, commitFile, verbosity)

				if len(Message) == 0 {
					Message = filepath.Base(commitFile)

					if verbosity >= 1 {
						fmt.Printf("Message not specified, setting to: %s\n", Message)
					}
				}
			}

			dupver.CommitFile(tarFile, Message, verbosity)
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

				if verbosity >= 1 {
					fmt.Printf("Message not specified, setting to: %s\n", Message)
				}				
			}

			tarFile := CreateTar(containingFolder, workdirFolder, verbosity)
			dupver.CommitFile(tarFile, Message, verbosity)			
		}
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	commitCmd.Flags().StringVarP(&Message, "message", "m", "", "Commit message")
}
