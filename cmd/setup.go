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
	Short: "Setup ian working environment",
	Long: `Ian requires you to be able to interact with Github through Git CLI.

    With projects subcommand being one of the core function of Ian, setup will
    install what is necessessary to deploy on GCE.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting setup")

		if _, err := os.Stat(OSPackageManager.GetExecPath()); err != nil {
			fmt.Println("Missing OS package manager !")
			os.Exit(1)
		}

		config.SetupConfigFiles()
		OSPackageManager.Setup()

		env.SetupDotFiles(
			config.Vipers["config"].GetString("github_username"),
			config.DotfilesDirPath,
		)

		packageManagerNames := config.Vipers["env"].AllKeys()
		for _, packageManagerName := range packageManagerNames {
			packages := config.Vipers["env"].GetStringSlice(packageManagerName)
			packageManager := pm.GetPackageManager(packageManagerName)

			env.SetupPackages(packageManager, packages)
		}
		fmt.Println("Ending setup.")
	},
}
