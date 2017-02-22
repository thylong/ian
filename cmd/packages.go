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
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/thylong/ian/backend/command"
	"github.com/thylong/ian/backend/config"
	"github.com/thylong/ian/backend/packages"

	"github.com/spf13/cobra"
)

// packagesCmd represents the packages command
var packagesCmd = &cobra.Command{
	Use:   "packages",
	Short: "packages allows you to manage ian extensions",
	Long: `packages allows you to install, uninstall subcommands for Ian.

    An example would be baily CLI (a nice bot powered by @samueldelesque & Dailymotion)`,
}

var selectedPackageManager string

func init() {
	installPackagesCmd.Flags().StringVarP(&selectedPackageManager, "package-manager", "p", "", "Package manager to use to install.")
	uninstallPackagesCmd.Flags().StringVarP(&selectedPackageManager, "package-manager", "p", "", "Package manager to use to install.")
	RootCmd.AddCommand(packagesCmd)
	packagesCmd.AddCommand(listPackagesCmd)
	packagesCmd.AddCommand(installPackagesCmd)
	packagesCmd.AddCommand(uninstallPackagesCmd)
	packagesCmd.AddCommand(searchPackagesCmd)
}

var listPackagesCmd = &cobra.Command{
	Use:   "list",
	Short: "list installed ian extensions",
	Long: `list installed ian extensions,

    This won't list npm, pip, gem, composer or other kind of packages.`,
	Run: func(cmd *cobra.Command, args []string) {
		packagesUsages := `Package Commands:`
		for packageName, packageMeta := range config.Config.Packages {
			packagesUsages += "\n" + `  ` + packageName + ` ` + packageMeta["description"] + ` type:` + packageMeta["type"]
		}
		fmt.Println(packagesUsages)
	},
}

var installPackagesCmd = &cobra.Command{
	Use:   "install",
	Short: "install a new extension",
	Long: `install new subcommand(s) for Ian.

    An example would be baily CLI (a nice bot powered by @samueldelesque & Dailymotion)`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			results, err := packages.IsAvailableOnPackageManagers(arg)
			if err != nil {
				log.Error(err.Error())
				return
			}

			if selectedPackageManager == "" {
				for packageManager, available := range results {
					if available {
						fmt.Println("Package found on :", packageManager)
					}
				}
				fmt.Println("\nUse -p option to install,\nian packages --help to print usage.")
			} else {
				cmdParams := []string{}
				installParams := strings.Split(config.Config.Managers[selectedPackageManager]["install_cmd"], " ")
				cmdParams = append(installParams, arg)

				termCmd := exec.Command(selectedPackageManager, cmdParams...)
				command.ExecuteCommand(termCmd)

				err := packages.WritePackageEntry(selectedPackageManager, arg)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	},
}

var searchPackagesCmd = &cobra.Command{
	Use:   "search",
	Short: "search a new extension",
	Long: `search subcommand(s) for Ian.

    An example would be baily CLI (a nice bot powered by @samueldelesque & Dailymotion)`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			results, err := packages.SearchOnPackageManagers(arg)
			if err != nil {
				log.Error(err.Error())
				return
			}
			fmt.Println(results)
		}
	},
}

var uninstallPackagesCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "uninstall an extension",
	Long:  `uninstall an extension.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			results, err := packages.IsAvailableOnPackageManagers(arg)
			if err != nil {
				log.Error(err.Error())
				return
			}

			if selectedPackageManager == "" {
				for packageManager, available := range results {
					if available {
						fmt.Println("Package found on :", packageManager)
					}
				}
				fmt.Println("\nUse -p option to uninstall,\nian packages --help to print usage.")
			} else {
				cmdParams := []string{}
				installParams := strings.Split(config.Config.Managers[selectedPackageManager]["uninstall_cmd"], " ")
				cmdParams = append(installParams, arg)
				fmt.Println(cmdParams)

				termCmd := exec.Command(selectedPackageManager, cmdParams...)
				command.ExecuteCommand(termCmd)

				err := packages.UnwritePackageEntry(selectedPackageManager, arg)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	},
}
