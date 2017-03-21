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
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/thylong/ian/backend/command"
	"github.com/thylong/ian/backend/config"
	"github.com/thylong/ian/backend/projects"
	"github.com/thylong/ian/backend/repo"
)

var customCmdDescription string
var softDeleteCmdParam bool

func init() {
	RootCmd.AddCommand(projectCmd)

	for pName, pCmd := range config.GetProjects() {
		projectCmd.AddCommand(pCmd)

		deleteProjectCmd := deleteProjectCmd()
		setProjectCmd := setProjectCmd()

		deleteProjectCmd.Flags().BoolVar(&softDeleteCmdParam, "soft", false, "Delete only ian configuration but keeps history")
		setProjectCmd.Flags().StringVarP(&customCmdDescription, "description", "d", "", "Description of the custom project command (40 characters maximum)")

		pCmd.AddCommand(
			statusProjectCmd(),
			statsProjectCmd(),
			configProjectCmd(),
			deployProjectCmd(),
			rollbackProjectCmd(),
			dbProjectCmd(),
			addProjectCmd(),
			deleteProjectCmd,
			unsetProjectCmd(),
			cloneCmd(),
			setProjectCmd,
		)

		for _, customCmd := range config.GetCustomCmds(pName) {
			pCmd.AddCommand(customCmd)
		}
	}
}

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Interact with local projects",
	Long:  `Interact with a project using predefined commands, or define custom commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(config.Vipers["projects"].AllSettings()) == 0 {
			fmt.Println("/!\\ You currently have no projects set up.")
			in := config.GetUserInput("Would you like to add one using 'ian project new'? (Y/n)")
			if strings.ToLower(in) != "y" && strings.ToLower(in) != "yes" && strings.ToLower(in) != "" {
				return
			}
			addProjectCmd().Execute()
		}
		cmd.Usage()
	},
}

func statusProjectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Display health status",
		Long:  `Display health status.`,
		Run: func(cmd *cobra.Command, args []string) {
			baseURL := config.Vipers["projects"].GetStringMapString(cmd.Parent().Use)["url"]
			healthEndpoint := config.Vipers["projects"].GetStringMapString(cmd.Parent().Use)["health"]
			projects.Status(cmd.Parent().Use, baseURL, healthEndpoint)
			repo.Status(cmd.Parent().Use)
		},
	}
}

func statsProjectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Display Github stat's for a project",
		Long:  `Display Github stat's for a project.`,
		Run: func(cmd *cobra.Command, args []string) {
			repositoryURL := fmt.Sprintf(
				"https://api.github.com/repos/%s",
				config.Vipers["projects"].GetStringMapString(cmd.Parent().Use)["repository"],
			)
			projects.Stats(cmd.Parent().Use, repositoryURL)
		},
	}
}

func configProjectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Display a project's configuration",
		Long:  `Display a project's configuration.`,
		Run: func(cmd *cobra.Command, args []string) {
			settings := config.Vipers["projects"].GetStringMap(cmd.Parent().Use)
			prettySettings, _ := json.MarshalIndent(settings, "", "  ")
			fmt.Printf("%s configuration:\n%s\n}", cmd.Parent().Use, prettySettings)
		},
	}
}

func deployProjectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a new version of a project",
		Long:  `Deploy a new version of a project.`,
		Run: func(cmd *cobra.Command, args []string) {
			deployCmdContent := config.Vipers["projects"].GetStringMapString(cmd.Parent().Use)["deploy_cmd"]
			deployCmd := strings.Split(deployCmdContent, " ")
			termCmd := exec.Command(deployCmd[0], deployCmd[1:]...)
			termCmd.Dir = config.Vipers["config"].GetString("repositories_path")
			command.ExecuteInteractiveCommand(termCmd)
		},
	}
}

func rollbackProjectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rollback",
		Short: "Rollback to a previous version of a project",
		Long:  `Rollback to a previous version of a project.`,
		Run: func(cmd *cobra.Command, args []string) {
			rollbackCmdContent := config.Vipers["projects"].GetStringMapString(cmd.Parent().Use)["rollback_cmd"]
			rollbackCmd := strings.Split(rollbackCmdContent, " ")
			termCmd := exec.Command(rollbackCmd[0], rollbackCmd[:1]...)
			termCmd.Dir = config.Vipers["config"].GetString("repositories_path")
			command.ExecuteInteractiveCommand(termCmd)
		},
	}
}

func dbProjectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "db",
		Short: "Connect to the project's database",
		Long:  `Connect to the project's database.`,
		Run: func(cmd *cobra.Command, args []string) {
			rollbackCmdContent := config.Vipers["projects"].GetStringMapString(cmd.Parent().Use)["db_cmd"]
			dbCmd := strings.Split(rollbackCmdContent, " ")
			termCmd := exec.Command(dbCmd[0], dbCmd[:1]...)
			command.ExecuteInteractiveCommand(termCmd)
		},
	}
}

func addProjectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "Add a new project configuration",
		Long:  `Add a new project configuration.`,
		Run: func(cmd *cobra.Command, args []string) {
			NewP := make(map[string]string)
			projectName := config.GetUserInput("Enter the name of the project")
			NewP["description"] = config.GetUserInput("Enter the project description")
			NewP["url"] = config.GetUserInput("Enter the project URL (example: http://goian.io)")
			NewP["repository"] = config.GetUserInput("Enter the project repository (example: thylong/ian)")
			NewP["health"] = config.GetUserInput("Enter the health check relative URL (example: /status)")
			NewP["db_cmd"] = config.GetUserInput("Enter the db connection command (example: mongo localhost)")
			NewP["deploy_cmd"] = config.GetUserInput("Enter the deploy command (example: bash deploy.sh)")
			NewP["rollback_cmd"] = config.GetUserInput("Enter the rollback command (example: bash rollback.sh)")

			projectsContent := config.Vipers["projects"].AllSettings()
			projectsContent[projectName] = NewP
			config.UpdateYamlFile(
				config.ConfigFilesPathes["projects"],
				projectsContent,
			)
		},
	}
}

func deleteProjectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Delete a project configuration",
		Long:  `Delete a project configuration.`,
		Run: func(cmd *cobra.Command, args []string) {
			if !softDeleteCmdParam {
				repo.Remove(cmd.Parent().Use)
			}

			editContent := config.Vipers["projects"].AllSettings()
			delete(editContent, cmd.Parent().Use)

			config.UpdateYamlFile(
				config.ConfigFilesPathes["projects"],
				config.Vipers["projects"].AllSettings(),
			)
		},
	}
}

func setProjectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set",
		Short: "Define a subcommand",
		Long:  `Define a subcommand.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				fmt.Fprintf(os.Stderr, "%v Not enough argument.\n\n", color.RedString("Error:"))
				cmd.Usage()
				os.Exit(1)
			}
			if 5 < len(customCmdDescription) && len(customCmdDescription) < 40 {
				fmt.Fprintf(os.Stderr, "%v Description must be between 5 and 40 alphanumeric characters.\n\n", color.RedString("Error:"))
				cmd.Usage()
				os.Exit(1)
			}

			editContent := config.Vipers["projects"].GetStringMapString(cmd.Parent().Use)
			editContent[fmt.Sprintf("%s_custom_cmd", args[0])] = fmt.Sprintf("%s=%s", customCmdDescription, strings.Join(args[1:], " "))
			config.Vipers["projects"].Set(cmd.Parent().Use, editContent)

			projectsContent := config.Vipers["projects"].AllSettings()
			config.UpdateYamlFile(
				config.ConfigFilesPathes["projects"],
				projectsContent,
			)
		},
	}
}

func cloneCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clone",
		Short: "Clone project repository",
		Long: `Clone the project repository.

        Clone the project repository in repositories_path path.`,
		Run: func(cmd *cobra.Command, args []string) {
			editContent := config.Vipers["projects"].GetStringMapString(cmd.Parent().Use)
			if projectRepository, ok := editContent["repository"]; ok {
				repo.Clone(projectRepository)
			}
		},
	}
}

func unsetProjectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unset",
		Short: "Remove a subcommand",
		Long:  `Remove a subcommand.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Fprintf(os.Stderr, "%v Not enough argument.\n\n", color.RedString("Error:"))
				cmd.Usage()
				os.Exit(1)
			}

			editContent := config.Vipers["projects"].GetStringMapString(cmd.Parent().Use)
			delete(editContent, fmt.Sprintf("%s_cmd", args[0]))

			config.UpdateYamlFile(
				config.ConfigFilesPathes["projects"],
				config.Vipers["projects"].AllSettings(),
			)
		},
	}
}
