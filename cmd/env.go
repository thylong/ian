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
	"time"

	"github.com/spf13/cobra"
	"github.com/thylong/ian/backend/config"
	"github.com/thylong/ian/backend/env"
)

func init() {
	RootCmd.AddCommand(envCmd)
	envCmd.AddCommand(envDescribeCmd)
	envCmd.AddCommand(envUpdateCmd)
	envCmd.AddCommand(envUpgradeCmd)
	envCmd.AddCommand(envSaveCmd)
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
		OSPackageManager.UpgradeAll()
	},
}

var envSaveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save current configuration files to the dotfiles repository",
	Long:  `Save current configuration files to the dotfiles repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		env.Save(
			config.DotfilesDirPath,
			config.Vipers["config"].GetString("dotfiles_repository"),
			config.Vipers["config"].GetString("default_save_message"),
			[]string{".testong"},
		)
	},
}
