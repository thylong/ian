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
	envCmd.AddCommand(envInfoCmd)
	envCmd.AddCommand(envUpdateCmd)
	envCmd.AddCommand(envSaveCmd)
}

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Get infos about the local environment",
	Long: `Get general or detailes informations about the current environment.
Currently implemented: System, Network, Security, current load.`,
}

var envInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get infos about the local environment",
	Long:  `Get general or detailes informations about the current environment. Currently implemented : System, Network, Security, current load.`,
	Run: func(cmd *cobra.Command, args []string) {
		env.GetInfos()
	},
}

var envUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the local environment",
	Long:  `Update the local environment with infos stored in the local config.`,
	Run: func(cmd *cobra.Command, args []string) {
		func() {
			go func() {
				for {
					for _, v := range `-\|/` {
						fmt.Printf("\r Updating env... %c", v)
						time.Sleep(100 * time.Millisecond)
					}
				}
			}()
			OSPackageManager.UpdateAll()
		}()

		OSPackageManager.UpgradeAll()
		OSPackageManager.Cleanup()
	},
}

var envSaveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save current configuration to distant dotfiles repositories",
	Long:  `Move current configuration files of the user to a dotfiles sub-directory (if not exists), create symlinks to previous place, then finally create and push the repositories on github`,
	Run: func(cmd *cobra.Command, args []string) {
		env.Save(
			config.DotfilesDirPath,
			config.Vipers["config"].GetString("dotfiles_repository"),
			config.Vipers["config"].GetString("default_save_message"),
			[]string{".testong"},
		)
	},
}
