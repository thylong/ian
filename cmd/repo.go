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
	"github.com/spf13/cobra"
	"github.com/thylong/ian/backend/repo"
)

func init() {
	RootCmd.AddCommand(repoCmd)
	// Subcommands
	repoCmd.AddCommand(cleanCmd)
}

// repoCmd represents the repo command
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage local repositories",
	Long:  `Manage local repositories.`,
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
