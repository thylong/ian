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
	"os"

	"github.com/thylong/ian/backend/config"
	pm "github.com/thylong/ian/backend/package-managers"

	"github.com/spf13/cobra"
)

// PackageManagerFlag contains the value of package-manager flag.
var PackageManagerFlag string

func init() {
	installPackagesCmd.Flags().StringVarP(&PackageManagerFlag, "package-manager", "p", "", "Package manager to use to install.")
	uninstallPackagesCmd.Flags().StringVarP(&PackageManagerFlag, "package-manager", "p", "", "Package manager to use to install.")
	RootCmd.AddCommand(packagesCmd)
	packagesCmd.AddCommand(listPackagesCmd)
	packagesCmd.AddCommand(installPackagesCmd)
	packagesCmd.AddCommand(uninstallPackagesCmd)
	packagesCmd.AddCommand(searchPackagesCmd)
}

// packagesCmd represents the packages command
var packagesCmd = &cobra.Command{
	Use:   "packages",
	Short: "packages allows you to manage ian extensions",
	Long: `packages allows you to install, uninstall subcommands for Ian.

    An example would be baily CLI (a nice bot powered by @samueldelesque & Dailymotion)`,
}

var listPackagesCmd = &cobra.Command{
	Use:   "list",
	Short: "list installed ian extensions",
	Long: `list installed ian extensions,

    This won't list npm, pip, gem, composer or other kind of packages.`,
	Run: func(cmd *cobra.Command, args []string) {
		packagesUsages := `Package Commands:`
		for packageName, packageMeta := range config.ConfigMap.Packages {
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
		packageManager := pm.GetPackageManager(PackageManagerFlag)

		for _, arg := range args {
			packageManager.Install(arg)
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
			_, err := pm.SearchOnPackageManagers(arg)
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				return
			}
		}
	},
}

var uninstallPackagesCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "uninstall an extension",
	Long:  `uninstall an extension.`,
	Run: func(cmd *cobra.Command, args []string) {
		packageManager := pm.GetPackageManager(PackageManagerFlag)
		for _, arg := range args {
			packageManager.Uninstall(arg)
		}
	},
}
