package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/restic/chunker"
	"github.com/spf13/cobra"

	"github.com/akbarnes/dupver/src/dupver"
	"github.com/akbarnes/dupver/src/fancyprint"
)

var UseDefaultRepo bool
var EditRepoCfg bool

// initCmd represents the init command
var repoConfigCmd = &cobra.Command{
	Use:   "config [path]",
	Short: "Initialize a repository",
	Long: `This initializes a dupver repository (as oppposed to a workdir)

If no arguments are provided this command will create a repository in
the current working directory. The first optional positional argument
allows for the repository path to be specified. The second optional
positional argument allows for the repository name to be specified.
if no repository name is specified, the repository takes on the default
name of "main."`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := dupver.Options{JsonOutput: JsonOutput}
		prefs, _ := dupver.ReadPrefs(opts)
		fancyprint.Setup(Debug, Verbose, Quiet, Monochrome)
		// repoName := RepoName
		repoPath := RepoPath

		if UseDefaultRepo {
			fancyprint.Debug("Using default repo")
			repoPath = prefs.DefaultRepo
		} else if len(args) >= 1 {
			repoPath = args[0]
		}

		// if len(args) >= 2 {
		// 	repoName = args[1]
		// }

		if len(repoPath) == 0 {
			prefs, _ := dupver.ReadPrefs(opts)
			repoPath = prefs.DefaultRepo
			fancyprint.Noticef("Using default repo: %s\n", repoPath)
		}

		if EditRepoCfg {
			cfgPath := filepath.Join(repoPath, "config.toml")
			dupver.EditFile(cfgPath, prefs)
			return
		}

		// dupver.InitRepo(repoPath, repoName, ChunkerPolynomial, CompressionLevel, opts)
		cfg := dupver.ReadRepoConfig(repoPath)

		if len(args) == 2 || (UseDefaultRepo && len(args) == 1) {
			key := args[0]

			if !UseDefaultRepo {
				key = args[1]
			}

			switch key {
			case "Version":
				dupver.MultiPrint(cfg.Version, opts)
			case "ChunkerPolynomial":
				dupver.MultiPrint(cfg.ChunkerPolynomial, opts)
			case "CompressionLevel":
				dupver.MultiPrint(cfg.CompressionLevel, opts)
			default:
				fancyprint.Warnf("Key %s doesn't exit in the repository configuration.", key)

				if JsonOutput {
					dupver.PrintJson(nil)
				}
			}

			return
		}

		if len(args) == 3 || (UseDefaultRepo && len(args) == 2) {
			key := args[0]
			val := args[1]

			if !UseDefaultRepo {
				key = args[1]
				val = args[2]
			}

			switch key {
			case "Version":
				v, err := strconv.ParseInt(val, 0, 32)
				dupver.Check(err)
				cfg.Version = int(v)
			case "ChunkerPolynomial":
				p, err := strconv.ParseInt(val, 0, 64)
				dupver.Check(err)
				cfg.ChunkerPolynomial = chunker.Pol(p)
			case "CompressionLevel":
				c, err := strconv.ParseUint(val, 0, 16)
				dupver.Check(err)
				cfg.CompressionLevel = uint16(c)
			default:
				fancyprint.Warnf("Key %s doesn't exit in the repository configuration.", key)

				if JsonOutput {
					dupver.PrintJson(nil)
				}

				return
			}

			dupver.SaveRepoConfig(repoPath, cfg, true)
		}

		if JsonOutput {
			dupver.PrintJson(cfg)
		} else {
			cfg.Print()
		}
	},
}

func init() {
	repoCmd.AddCommand(repoConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")
	repoConfigCmd.Flags().BoolVarP(&UseDefaultRepo, "default", "D", false, "use default repo")
	repoConfigCmd.Flags().BoolVarP(&EditRepoCfg, "edit", "e", false, "edit the repo config with the specified editor")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
