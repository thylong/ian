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
	"strings"

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
		fmt.Print([]byte(`Welcome to Ian!
Ian is a simple tool to manage your development environment, repositories,
and projects.

Learn more about Ian at http://goian.io

To benefit from all of Ian’s feature’s, you’ll need to provide:
    - A working OS Package Manager (will set up if missing)
    - The full path of your repositories (example: /Users/thylong/repositories)
    - The path of your dotfiles Github repository (example: thylong/dotfiles)

`))
		setupBool := config.GetUserInput("Do you want to set up Ian now? (Y/n): ")
		if strings.ToLower(setupBool) != "y" && strings.ToLower(setupBool) != "yes" && strings.ToLower(setupBool) != "" {
			fmt.Println("You're ready to start using Ian. Note that if you try to use some of Ian's\nfeatures you'll be prompted for these details again.")
			os.Exit(1)
		}
		if _, err := os.Stat(OSPackageManager.GetExecPath()); err != nil {
			fmt.Println("Installing OS package manager...")
			err = OSPackageManager.Setup()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Missing OS package manager !")
				os.Exit(1)
			}
		}

		env.SetupDotFiles(
			config.Vipers["config"].GetString("dotfiles_repository"),
			config.DotfilesDirPath,
		)

		packageManagerNames := config.Vipers["env"].AllKeys()
		for _, packageManagerName := range packageManagerNames {
			packages := config.Vipers["env"].GetStringSlice(packageManagerName)
			packageManager := pm.GetPackageManager(packageManagerName)

			env.SetupPackages(packageManager, packages)
		}
		fmt.Println("Great! You're ready to start using Ian.")
	},
}
