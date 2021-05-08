package cmd

import (
	// "fmt"
	// "log"
	// "path"
	"os"
	// "os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/akbarnes/dupver/src/dupver"	
	"github.com/akbarnes/dupver/src/fancyprint"
)

var Add bool
var Message string


// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit [tar_file]",
	Short: "Commit the current working directory or a tarball",
	Long: `This commits the current working directory or a tarball.

	If no arguments are provided this command will commit the
	current working directory (if initialized as a dupver
	working directory). If a single positional argument is
	provided then a tarball with the specified file name
	is committed. This is intended to allow for git-style
	incremental commits using the -r option of tar. The 
	commit command does not require a commit message, though
	this can be specified with the --message flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := dupver.Options{}
		fancyprint.Setup(Debug, Verbose, Quiet, Monochrome)

        if len(WorkDirPath) > 0 {
            os.Chdir(WorkDirPath)
        }

		fancyprint.Setup(Debug, Verbose, Quiet, Monochrome)
		cfg, err := dupver.ReadWorkDirConfig(WorkDirPath)

		if err != nil {
			// Todo: handle invalid configuration file
			fancyprint.Warn("Could not read configuration file. Has the project working directory been initialized?")
			os.Exit(1)
		}

		if len(WorkDirPath) > 0 {
			os.Chdir(WorkDirPath)
		}

		fancyprint.Debugf("Workdir Configuration: %+v\n", cfg)
		fancyprint.Debugf("Repo name: %s\nRepo path: %s\n", RepoName, RepoPath)

		if len(RepoName) == 0 {
			RepoName = cfg.DefaultRepo
		}

		if len(RepoPath) == 0 {
			RepoPath = cfg.Repos[RepoName]
			fancyprint.Debugf("Updating repo path to: %s\n", RepoPath)
		}

		if len(Branch) == 0 {
			Branch = cfg.Branch
		}		


		opts.RepoName = RepoName
		opts.RepoPath = RepoPath
		opts.Branch = Branch

		workDir := dupver.InstantiateWorkDir(cfg)

		if AllBranches {
			opts.Branch = ""
		}

		if len(args) >= 1 {
			commitFile := args[0]
			tarFile := commitFile
			containingFolder := filepath.Dir(commitFile)

			if !strings.HasSuffix(commitFile, "tar") {
				fancyprint.Debugf("%s -> %s, %s\n", commitFile, containingFolder, commitFile)
				tarFile = dupver.CreateTar(containingFolder, commitFile)

				if len(Message) == 0 {
					Message = filepath.Base(commitFile)
					fancyprint.Debugf("Message not specified, setting to: %s\n", Message)
				}
			}

			parentIds := []string{}
			workDir.CommitFile(tarFile, parentIds, Message, JsonOutput)
		} else {
			workDir.Commit(Message, JsonOutput)
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
	commitCmd.Flags().BoolVarP(&Add, "add", "a", false, "Unused, but added for git compatibility")
}
