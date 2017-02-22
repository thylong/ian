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
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thylong/ian/backend/command"
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
		for _, arg := range args {
			projectConfig := viper.GetStringMap("projects")[arg]
			baseURL := projectConfig.(map[interface{}]interface{})["url"]
			healthEndpoint := projectConfig.(map[interface{}]interface{})["health"]

			url := baseURL.(string) + healthEndpoint.(string)
			resp, err := http.Get(url)
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
			}
			defer resp.Body.Close()

			if statusCode := resp.StatusCode; statusCode == 200 {
				fmt.Printf("%s : OK", arg)
			} else {
				fmt.Printf("%s : ERROR", arg)
			}
		}
	},
}

var statsProjectCmd = &cobra.Command{
	Use:   "stats",
	Short: "Get the number of stars and forks for projects",
	Long:  `Get the number of stars and forks for projects.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			projectConfig := viper.GetStringMap("projects")[arg]
			repositoryURL := projectConfig.(map[interface{}]interface{})["repository"]
			url := "https://api.github.com/repos/" + repositoryURL.(string)

			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("Error : %s", err.Error())
			}
			content, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()

			var jsonContent map[string]interface{}
			err = json.Unmarshal(content, &jsonContent)
			if err != nil {
				fmt.Printf("Error : %s", err.Error())
				return
			}

			fmt.Print(arg + " :")
			fmt.Printf("- Forks: %v", jsonContent["forks_count"])
			fmt.Printf("- Stars: %v", jsonContent["stargazers_count"])
			fmt.Printf("- Open Issues: %v", jsonContent["open_issues_count"])
			fmt.Printf("- Last update: %v", jsonContent["updated_at"])
		}
	},
}

var configProjectCmd = &cobra.Command{
	Use:   "config",
	Short: "Gather general config of projects",
	Long:  `Gather general config of projects.`,
	Run: func(cmd *cobra.Command, args []string) {
		projectConfig := viper.GetStringMap("projects")
		for _, project := range projectConfig {
			fmt.Print(project)
		}
	},
}

var deployProjectCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a new version of a project",
	Long:  `Deploy a new version of a project.`,
	Run: func(cmd *cobra.Command, args []string) {
		projectConfig := viper.GetStringMap("projects")[args[0]]
		deployCmd := strings.Split(projectConfig.(map[interface{}]interface{})["deploy_cmd"].(string), " ")
		termCmd := exec.Command(deployCmd[0], deployCmd[:1]...)
		termCmd.Dir = viper.GetString("repositories_path")
		command.ExecuteCommand(termCmd)
	},
}

var rollbackProjectCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback to a previous version of a project",
	Long:  `Rollback a previous version of a project.`,
	Run: func(cmd *cobra.Command, args []string) {
		projectConfig := viper.GetStringMap("projects")[args[0]]
		rollbackCmd := strings.Split(projectConfig.(map[interface{}]interface{})["rollback_cmd"].(string), " ")
		termCmd := exec.Command(rollbackCmd[0], rollbackCmd[:1]...)
		termCmd.Dir = viper.GetString("repositories_path")
		command.ExecuteCommand(termCmd)
	},
}

var dbProjectCmd = &cobra.Command{
	Use:   "db",
	Short: "Connect to the Database of the project",
	Long:  `Connect to the Database of the project.`,
	Run: func(cmd *cobra.Command, args []string) {
		projectConfig := viper.GetStringMap("projects")[args[0]]
		dbCmd := strings.Split(projectConfig.(map[interface{}]interface{})["db_cmd"].(string), " ")
		termCmd := exec.Command(dbCmd[0], dbCmd[:1]...)
		termCmd.Stdout = os.Stdout
		termCmd.Stdin = os.Stdin
		termCmd.Stderr = os.Stderr
		termCmd.Run()
	},
}
