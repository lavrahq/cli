// Copyright © 2019 Scott Plunkett <plunkets@aeoss.io>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	"github.com/lavrahq/cli/util/logs"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "runctl",
	Short: "Allows deploying, configuring, and extensing Run stacks.",
	Long: `This CLI offers an easy way to deploy, configure, and extend Run
stacks, including the open source and freely available Lavra products. We provide
this CLI utility to handle the administrative tasks so that administrators
do not need advanced knowledge of orchestration platforms in order to securely
run stacks.`,
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lavra/config.yml)")

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
		home, err := homedir.Expand("~/.lavra")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		file, _ := homedir.Expand("~/.lavra/config.yml")
		if _, err := os.Stat(file); os.IsNotExist(err) {
			err := os.MkdirAll(home, os.ModePerm)
			if err != nil {
				fmt.Println("ERROR: Could not create .lavra config dir. Config will not load.")

				return
			}

			_, err = os.Create(file)
			if err != nil {
				fmt.Println("ERROR: Could not create config.yml file. Config will not load.")

				return
			}
		}

		// Search config in .lavra within home directory with name "config" (without extension).
		viper.SetConfigName("config")

		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
	}

	logs.InitGlobalLogging()

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Failed to read in config file. Please check the config.yml file.")

		return
	}
}
