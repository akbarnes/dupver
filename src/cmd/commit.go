package cmd

import (
	"fmt"
	"log"
	// "path"
	"path/filepath"
	"strings"
	"os"
	"os/exec"

	"github.com/akbarnes/dupver/src/dupver"
	"github.com/spf13/cobra"
)

var Add bool
var Message string
var ParentCommitIds string

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
	this can be specified with the --message flag. The
	--parent flag bears mentioning. Dupver does not currently
	support merge functionality, though it tracks branches.
	By default a commit's parent will be the commit id
	of the current branch head. The --parent flag allows for
	a list of commits to be specified instead which is used
	after performing a merge. Merges must be performed manually 
	by checking out the branches into separate project working
	directories and using an external diff/merge tool to resolve 
	the differences. Once this is performed, the --parent
	flag is used to keep track of this merge.`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := dupver.SetVerbosity(dupver.Options{Color: true}, Debug, Verbose, Quiet)
		
		if Monochrome || Quiet {
			opts.Color = false
		}

		opts.RepoName = RepoName
		opts.RepoPath = RepoPath
		opts.BranchName = BranchName
		
		parentIds := []string{}
		unfilteredParentIds := strings.Split(ParentCommitIds, ",")

		for i := range unfilteredParentIds {
			if len(unfilteredParentIds[i]) > 0 {
				parentIds = append(parentIds, unfilteredParentIds[i])
			}
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

			dupver.CommitFile(tarFile, parentIds, Message, opts)
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
			dupver.CommitFile(tarFile, parentIds, Message, opts)	
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
	commitCmd.Flags().StringVarP(&ParentCommitIds, "parent", "c", "", "Comma separated list of parent commit ID(s)")
}
