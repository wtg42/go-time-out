/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v6"
	"github.com/vbauerster/mpb/v6/decor"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var breaktime string
var worktime string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "time-out",
	Short: "Take a break",
	Long: `Remind you pendding current job and take a break.
	use time-out command to set working time and break time`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Break for?(seconds): ")
		scanner.Scan()
		breaktime = scanner.Text()

		fmt.Print("Every?(minutes): ")
		scanner.Scan()
		worktime = scanner.Text()
		seconds, _ := strconv.Atoi(breaktime)

		timer := time.NewTicker(time.Second * 1)

		// initialize progress container, with custom width
		p := mpb.New(
			mpb.WithWidth(60),
			mpb.WithRefreshRate(10*time.Millisecond),
		)
		total := seconds
		name := "Working time:"
		// frames := []string{"/", "\\"}
		// adding a single bar, which will inherit container's width
		bar := p.Add(int64(total),
			// progress bar filler with customized style
			mpb.NewBarFiller("[=>-]"),
			mpb.PrependDecorators(
				// display our name with one space on the right
				decor.Name(name, decor.WC{W: len(name) + 10, C: decor.DidentRight}),
				// replace ETA decorator with "done" message, OnComplete event
				decor.OnComplete(
					decor.Counters(0, "% d / % d"), "done",
				),
			),
			mpb.AppendDecorators(decor.Spinner(nil)),
		)
		// start := time.Now()
		for {
			<-timer.C
			bar.Increment()
			if bar.Completed() {
				break
			}
		}
		p.Wait()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.demo.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVarP(&breaktime, "break", "b", "", "The seconds you want to break for.")
	rootCmd.Flags().StringVarP(&worktime, "work", "w", "", "Set a the duration of work.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".demo" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".demo")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
