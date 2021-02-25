package cmd

import (
	"fmt"
	"log"

	// "path"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/akbarnes/dupver/src/dupver"
	"github.com/spf13/cobra"
)

var Add bool
var Message string

// Write a project working directory to a tar file in temp
// given a working directory path and the path of its parent folder
func CreateTar(parentPath string, commitPath string, opts dupver.Options) string {
	tarFile := dupver.RandHexString(40) + ".tar"
	tarFolder := filepath.Join(dupver.GetHome(), "temp")
	tarPath := filepath.Join(tarFolder, tarFile)

	// InitRepo(workDir)
	if opts.Verbosity >= 1 {
		fmt.Printf("Tar path: %s\n", tarPath)
		fmt.Printf("Creating folder %s\n", tarFolder)
	}

	os.Mkdir(tarFolder, 0777)

	CompressTar(parentPath, commitPath, tarPath, opts)
	return tarPath
}

// Write a project working directory to a tar file
// given a working directory path, parent folder path and tar file path
func CompressTar(parentPath string, commitPath string, tarPath string, opts dupver.Options) string {
	if len(tarPath) == 0 {
		tarPath = commitPath + ".tar"
	}

	cleanCommitPath := filepath.Clean(commitPath)

	tarCmd := exec.Command("tar", "cfv", tarPath, cleanCommitPath)
	tarCmd.Dir = parentPath

	if opts.Verbosity >= 1 {
		log.Printf("Running tar cfv %s %s", tarPath, cleanCommitPath)
	}

	output, err := tarCmd.CombinedOutput()

	if err != nil {
		log.Fatal(fmt.Sprintf("Tar command failed\nOutput:\n%s\nError:\n%s\n", output, err))
	} else if opts.Verbosity >= 3 {
		fmt.Printf("Ran tar command with output:\n%s\n", output)
	}

	return tarPath
}

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
		opts := dupver.SetVerbosity(dupver.Options{Color: true}, Debug, Verbose, Quiet)

		if Monochrome || Quiet {
			opts.Color = false
		}

        if len(WorkDirPath) > 0 {
            os.Chdir(WorkDirPath)
        }

		opts.RepoName = RepoName
		opts.RepoPath = RepoPath
		opts.Branch = Branch

		if AllBranches {
			opts.Branch = ""
		}

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

			dupver.CommitFile(tarFile, []string{}, Message, opts)
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
			dupver.CommitFile(tarFile, []string{}, Message, opts)
			os.Remove(tarFile) // Delete the temporary file 
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
