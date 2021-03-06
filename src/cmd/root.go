package cmd

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var RepoName string
var RepoPath string
var Branch string
var AllBranches bool
var WorkDirPath string
var Debug bool
var Verbose bool
var Quiet bool
var Monochrome bool
var JsonOutput bool
var Color bool
var CompressionLevel uint16

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "main",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.main.yaml)")
	rootCmd.PersistentFlags().StringVarP(&RepoName, "repo-name", "n", "", "repository name")
	rootCmd.PersistentFlags().StringVarP(&RepoPath, "repo-path", "r", "", "repository path")
	rootCmd.PersistentFlags().StringVarP(&WorkDirPath, "workdir-path", "w", "", "project working directory path")
	rootCmd.PersistentFlags().StringVarP(&Branch, "branch", "b", "", "branch name")
	rootCmd.PersistentFlags().BoolVarP(&AllBranches, "all-branches", "A", false, "branch name")
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "debug output")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&Quiet, "quiet", "q", false, "quiet output")
	rootCmd.PersistentFlags().BoolVarP(&Monochrome, "monochrome", "M", false, "monochrome output")
	rootCmd.PersistentFlags().BoolVarP(&JsonOutput, "json", "j", false, "JSON output")
	rootCmd.PersistentFlags().Uint16VarP(&CompressionLevel, "compression-level", "C", zip.Deflate, "compression level")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".main" (without extension).
		viper.AddConfigPath(filepath.Join(home, ".dupver"))
		viper.SetConfigName("global_config.toml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
