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
	"github.com/thylong/ian/backend/repo"
)

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

// repoCmd represents the repo command
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage local repositories",
	Long:  `Manage local repositories.`,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List stored repositories",
	Long: `Manage local repositories.
It's currently possible to update remote, reset current branch to master, clean,
and delete.`,
	Run: func(cmd *cobra.Command, args []string) {
		repo.List()
	},
}

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone new repositories",
	Long: `Clone new repositories.

    Clone new repositories in repositories_path path.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			repo.Clone(arg)
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
			repo.Clean(arg)
		}
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update given repositories",
	Long: `Update given repositories.

    Update stored repositories in the repositories_path path.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			repo.UpdateAll()
		}
		for _, arg := range args {
			repo.UpdateOne(arg)
		}
	},
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade given repositories",
	Long: `Upgrade given repositories.

    Upgrade stored repositories in the repositories_path path.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			repo.UpgradeAll()
		}
		for _, arg := range args {
			repo.UpgradeOne(arg)
		}
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove given repositories",
	Long: `Remove given repositories.

    Remove stored repositories in the repositories_path path.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Fprint(os.Stderr, "At least one repository path should be provided.")
		}
		for _, arg := range args {
			repo.Remove(arg)
		}
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Status of given repositories",
	Long: `Status of given repositories.

    Give status of stored repositories in the repositories_path path.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Fprint(os.Stderr, "At least one repository path should be provided.")
		}
		for _, arg := range args {
			repo.Status(arg)
		}
	},
}
