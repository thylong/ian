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
	"strings"

	"github.com/spf13/cobra"
	"github.com/thylong/ian/backend/command"
	"github.com/thylong/ian/backend/config"
	"github.com/thylong/ian/backend/projects"
)

func init() {
	RootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(statusProjectCmd)
	projectCmd.AddCommand(statsProjectCmd)
	projectCmd.AddCommand(configProjectCmd)
	projectCmd.AddCommand(deployProjectCmd)
	projectCmd.AddCommand(rollbackProjectCmd)
	projectCmd.AddCommand(dbProjectCmd)
}

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Interact with listed project",
	Long:  `Get health statuses, general and detailed config about the listed projects.`,
}

var statusProjectCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the health statuses of projects",
	Long:  `Get the health statuses of the projects hosted versions.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, project := range args {
			baseURL := config.Vipers["projects"].GetStringMapString(project)["url"]
			healthEndpoint := config.Vipers["projects"].GetStringMapString(project)["health"]
			projects.Status(project, baseURL, healthEndpoint)
		}
	},
}

var statsProjectCmd = &cobra.Command{
	Use:   "stats",
	Short: "Get the number of stars and forks for projects",
	Long:  `Get the number of stars and forks for projects.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, project := range args {
			repositoryURL := fmt.Sprintf(
				"https://api.github.com/repos/%s",
				config.Vipers["projects"].GetStringMapString(project)["repository"],
			)
			projects.Stats(project, repositoryURL)
		}
	},
}

var configProjectCmd = &cobra.Command{
	Use:   "config",
	Short: "Gather general config of projects",
	Long:  `Gather general config of projects.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.Vipers["projects"].AllSettings())
	},
}

var deployProjectCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a new version of a project",
	Long:  `Deploy a new version of a project.`,
	Run: func(cmd *cobra.Command, args []string) {
		deployCmdContent := config.Vipers["projects"].GetStringMapString(args[0])["deploy_cmd"]
		deployCmd := strings.Split(deployCmdContent, " ")
		termCmd := exec.Command(deployCmd[0], deployCmd[1:]...)
		termCmd.Dir = config.Vipers["config"].GetString("repositories_path")
		command.ExecuteCommand(termCmd)
	},
}

var rollbackProjectCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback to a previous version of a project",
	Long:  `Rollback a previous version of a project.`,
	Run: func(cmd *cobra.Command, args []string) {
		rollbackCmdContent := config.Vipers["projects"].GetStringMapString(args[0])["deploy_cmd"]
		rollbackCmd := strings.Split(rollbackCmdContent, " ")
		termCmd := exec.Command(rollbackCmd[0], rollbackCmd[:1]...)
		termCmd.Dir = config.Vipers["config"].GetString("repositories_path")
		command.ExecuteCommand(termCmd)
	},
}

var dbProjectCmd = &cobra.Command{
	Use:   "db",
	Short: "Connect to the Database of the project",
	Long:  `Connect to the Database of the project.`,
	Run: func(cmd *cobra.Command, args []string) {
		rollbackCmdContent := config.Vipers["projects"].GetStringMapString(args[0])["db_cmd"]
		dbCmd := strings.Split(rollbackCmdContent, " ")
		termCmd := exec.Command(dbCmd[0], dbCmd[:1]...)
		termCmd.Stdout = os.Stdout
		termCmd.Stdin = os.Stdin
		termCmd.Stderr = os.Stderr
		termCmd.Run()
	},
}
