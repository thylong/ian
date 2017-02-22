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
	"os/exec"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thylong/ian/backend/command"
	"github.com/thylong/ian/backend/config"
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
		fmt.Println("====================")

		if _, err := os.Stat(OSPackageManager.GetExecPath()); err != nil {
			log.Fatal("Missing OS package manager !")
		}

		config.SetupConfigFile()
		config.SetupDotFiles()

		OSPackageManager.Setup()

		setupCLIPackages()
		setupGUIPackages()

		fmt.Println("====================")
		fmt.Println("Ending setup.")
	},
}

func setupCLIPackages() {
	fmt.Println("Installing CLI packages...")
	CLIPackages := viper.GetStringMapStringSlice("setup")["cli_packages"]

	if len(CLIPackages) == 0 {
		fmt.Println("No brew packages to install")
		return
	}

	for _, CLIPackage := range CLIPackages {
		OSPackageManager.Install(CLIPackage)
	}
}

func setupGUIPackages() {
	fmt.Println("Installing GUI packages...")
	GUIPackages := viper.GetStringMapStringSlice("setup")["gui_packages"]

	if len(GUIPackages) == 0 {
		fmt.Println("No GUI packages to install")
		return
	}

	for _, GUIPackage := range GUIPackages {
		command.ExecuteCommand(exec.Command("/usr/local/bin/brew", "cask", "install", GUIPackage))
	}
}
