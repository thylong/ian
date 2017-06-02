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

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/thylong/ian/backend/config"
	"github.com/thylong/ian/backend/env"
	pm "github.com/thylong/ian/backend/package-managers"
)

func init() {
	RootCmd.AddCommand(setupCmd)
}

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up ian configuration",
	Long: `Ian requires you to be able to interact with Github through Git CLI.

    With projects subcommand being one of the core function of Ian, setup will
    install what is necessessary to deploy on GCE.`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(OSPackageManager.GetExecPath()); err != nil {
			fmt.Println("Installing OS package manager...")
			if err = OSPackageManager.Setup(); err != nil {
				fmt.Fprintf(os.Stderr, "%v Missing OS package manager !", color.RedString("Error:"))
				return
			}
		}

		env.SetupDotFiles(
			config.Vipers["config"].GetStringMapString("dotfiles")["repository"],
			config.DotfilesDirPath,
		)
		// Refresh the configuration in case the imported dotfiels contains ian configuration
		config.RefreshVipers()

		fmt.Printf("You don't have any packages to be installed in your current ian configuration.\n")
		if _, ok := config.Vipers["env"]; !ok && config.GetBoolUserInput("Would you like to use a preset? (Y/n)") {
			in := config.GetUserInput(`Which preset would you like to use:
    1) Software engineer (generalist preset)
    2) Backend developer
    3) Frontend developer
    4) Ops
Enter your choice`)
			config.CreateEnvFileWithPreset(in)
		}

		packageManagers := config.Vipers["env"].AllKeys()
		for _, packageManager := range packageManagers {
			packages := config.Vipers["env"].GetStringSlice(packageManager)
			env.InstallPackages(pm.GetPackageManager(packageManager), packages)
		}
		fmt.Println("Great! You're ready to start using Ian.")
	},
}
