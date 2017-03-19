// Copyright © 2016 Theotime LEVEQUE theotime@protonmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/thylong/ian/backend/config"
	"github.com/thylong/ian/backend/env"
	pm "github.com/thylong/ian/backend/package-managers"
)

var allCmdParam bool

func init() {
	RootCmd.AddCommand(envCmd)

	envUpdateCmd.Flags().BoolVarP(&allCmdParam, "all", "a", false, "Run update on all Package managers")
	envUpgradeCmd.Flags().BoolVarP(&allCmdParam, "all", "a", false, "Run upgrade on all Package managers")
	envCmd.AddCommand(
		envDescribeCmd,
		envUpdateCmd,
		envUpgradeCmd,
		envSaveCmd,
	)
}

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage development environment",
	Long:  `Show details, update and save your development environment.`,
}

var envDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Show details of the current development environment",
	Long:  `Show details of the current development environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		env.Describe()
	},
}

var envUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the development environment",
	Long:  `Update the development environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		go func() {
			for {
				for _, v := range `-\|/` {
					fmt.Printf("\r Updating env... %c", v)
					time.Sleep(100 * time.Millisecond)
				}
			}
		}()
		if allCmdParam {
			for _, packageManager := range pm.SupportedPackageManagers {
				if packageManager.IsInstalled() {
					packageManager.UpdateAll()
				}
			}
			return
		}
		OSPackageManager.UpdateAll()
	},
}

var envUpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the development environment",
	Long:  `Upgrade the development environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		go func() {
			for {
				for _, v := range `-\|/` {
					fmt.Printf("\r Upgrading env... %c", v)
					time.Sleep(100 * time.Millisecond)
				}
			}
		}()
		if allCmdParam {
			for _, packageManager := range pm.SupportedPackageManagers {
				if packageManager.IsInstalled() {
					packageManager.UpgradeAll()
				}
			}
			return
		}
		OSPackageManager.UpgradeAll()
	},
}

var envSaveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save current configuration files to the dotfiles repository",
	Long:  `Save current configuration files to the dotfiles repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(config.Vipers["projects"].AllSettings()) == 0 {
			fmt.Println("/!\\ You currently have no defined path to your parent repositories directory.")
			in := config.GetUserInput("Would you like to provide the repositories_path now? (Y/n)")
			if strings.ToLower(in) != "y" && strings.ToLower(in) != "yes" && strings.ToLower(in) != "" {
				return
			}
		}
		env.Save(
			config.DotfilesDirPath,
			config.Vipers["config"].GetString("dotfiles_repository"),
			config.Vipers["config"].GetString("default_save_message"),
			[]string{".testong"},
		)
	},
}
