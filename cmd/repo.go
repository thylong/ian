// Copyright © 2016 THEOTIME LEVEQUE <theotime.leveque@gmail.com>
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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// repoCmd represents the repo command
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage stored repositories",
	Long: `Manage stored repositories.

    It's currently possible to update remote, reset current branch to master, clean, delete.
    /!\ The repositories_path must be set to the path of your repositories.`,
}

func init() {
	RootCmd.AddCommand(repoCmd)
	// Subcommands
	repoCmd.AddCommand(cleanCmd)
	repoCmd.AddCommand(cloneCmd)
	repoCmd.AddCommand(listCmd)
	repoCmd.AddCommand(removeCmd)
	repoCmd.AddCommand(statusCmd)
	repoCmd.AddCommand(updateCmd)
	repoCmd.AddCommand(upgradeCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List stored repositories",
	Long: `List stored repositories.

    List all stored repositories in repositories_path path.`,
	Run: func(cmd *cobra.Command, args []string) {
		termCmd := exec.Command("ls", "-lahs")
		fmt.Println(viper.GetString("repositories_path"))
		termCmd.Dir = viper.GetString("repositories_path")

		printFromCmdStds(termCmd)
	},
}

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone new repositories",
	Long: `Clone new repositories.

    Clone new repositories in repositories_path path.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			termCmd := exec.Command("git", "clone", "-v", arg)
			termCmd.Dir = viper.GetString("repositories_path")

			printFromCmdStds(termCmd)
		}
	},
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean given repositories",
	Long: `Clean given repositories.

    Clean stored repositories in the repositories_path path.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			termCmd := exec.Command("git", "clean", "-dffx", arg)
			termCmd.Dir = viper.GetString("repositories_path")

			printFromCmdStds(termCmd)
		}
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update given repositories",
	Long: `Update given repositories.

    Update stored repositories in the repositories_path path.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			termCmd := exec.Command("git", "fetch", arg)
			termCmd.Dir = viper.GetString("repositories_path")

			printFromCmdStds(termCmd)
		}
	},
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade given repositories",
	Long: `Upgrade given repositories.

    Upgrade stored repositories in the repositories_path path.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			termCmd := exec.Command("git", "pull", "--rebase", arg)
			termCmd.Dir = viper.GetString("repositories_path")

			printFromCmdStds(termCmd)
		}
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove given repositories",
	Long: `Remove given repositories.

    Remove stored repositories in the repositories_path path.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			termCmd := exec.Command("rm", "-rf", arg)
			termCmd.Dir = viper.GetString("repositories_path")

			printFromCmdStds(termCmd)
		}
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Status of given repositories",
	Long: `Status of given repositories.

    Give status of stored repositories in the repositories_path path.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			termCmd := exec.Command("git", "status")
			termCmd.Dir = viper.GetString("repositories_path") + "/" + arg

			printFromCmdStds(termCmd)
		}
	},
}
