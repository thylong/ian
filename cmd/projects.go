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
	"encoding/json"
	"fmt"
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
	projectCmd.AddCommand(addProjectCmd)
	projectCmd.AddCommand(deleteProjectCmd)
	projectCmd.AddCommand(setProjectCmd)
	projectCmd.AddCommand(unsetProjectCmd)
}

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Interact with local projects",
	Long:  `Interact with a project using predefined commands, or define custom commands.`,
}

var statusProjectCmd = &cobra.Command{
	Use:   "status",
	Short: "Display health status",
	Long:  `Display health status.`,
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
	Short: "Display Github stat's for a project",
	Long:  `Display Github stat's for a project.`,
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
	Short: "Display a project's configuration",
	Long:  `Display a project's configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, project := range args {
			settings := config.Vipers["projects"].GetStringMap(project)
			prettySettings, _ := json.MarshalIndent(settings, "", "  ")
			fmt.Printf("%s configuration:\n%s\n}", project, prettySettings)
		}
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
		command.ExecuteInteractiveCommand(termCmd)
	},
}

var rollbackProjectCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback to a previous version of a project",
	Long:  `Rollback to a previous version of a project.`,
	Run: func(cmd *cobra.Command, args []string) {
		rollbackCmdContent := config.Vipers["projects"].GetStringMapString(args[0])["deploy_cmd"]
		rollbackCmd := strings.Split(rollbackCmdContent, " ")
		termCmd := exec.Command(rollbackCmd[0], rollbackCmd[:1]...)
		termCmd.Dir = config.Vipers["config"].GetString("repositories_path")
		command.ExecuteInteractiveCommand(termCmd)
	},
}

var dbProjectCmd = &cobra.Command{
	Use:   "db",
	Short: "Connect to the project's database",
	Long:  `Connect to the project's database.`,
	Run: func(cmd *cobra.Command, args []string) {
		rollbackCmdContent := config.Vipers["projects"].GetStringMapString(args[0])["db_cmd"]
		dbCmd := strings.Split(rollbackCmdContent, " ")
		termCmd := exec.Command(dbCmd[0], dbCmd[:1]...)
		command.ExecuteInteractiveCommand(termCmd)
	},
}

var addProjectCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new project configuration",
	Long:  `Add a new project configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var deleteProjectCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a project configuration",
	Long:  `Delete a project configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var setProjectCmd = &cobra.Command{
	Use:   "set",
	Short: "Define a subcommand",
	Long:  `Define a subcommand.`,
	Run: func(cmd *cobra.Command, args []string) {
		// if len(args) < 2 {
		// 	return
		// }
		// projectContent := config.Vipers["projects"].GetStringMapString(args[0])
		// projectContent["set_cmd"] = strings.Join(args[1:], " ")
		//
		// out, err := yaml.Marshal(&projectContent)
		// if err != nil {
		// 	return
		// }
	},
}

var unsetProjectCmd = &cobra.Command{
	Use:   "unset",
	Short: "Remove a subcommand",
	Long:  `Remove a subcommand.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}
