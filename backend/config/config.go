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

package config

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thylong/ian/backend/command"
	"github.com/thylong/ian/backend/env"
)

// ConfigDirPath represents the path to config directory.
var ConfigDirPath string

// IanConfigPath represents the path to ian config directory.
var IanConfigPath string

// DotfilesDirPath represents the full path to the dotfiles directory.
var DotfilesDirPath string

// ConfigFilesPathes contains every config file pathes per filename.
var ConfigFilesPathes map[string]string

// Vipers contains all the initialized Vipers (config, env, projects)
var Vipers map[string]*viper.Viper

// YamlConfigMap is used to marshal/unmarshal config file.
type YamlConfigMap struct {
	Managers     map[string]map[string]string `json:"managers"`
	Repositories map[string]string            `json:"repositories"`
	Projects     map[string]map[string]string `json:"projects"`
	Setup        map[string][]string          `json:"setup"`
	Packages     map[string]map[string]string `json:"packages"`
}

// ConfigMap contains the config content.
var ConfigMap YamlConfigMap

func init() {
	usr, err := user.Current()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	ConfigDirPath = usr.HomeDir + "/.config"
	IanConfigPath = ConfigDirPath + "/ian/"
	DotfilesDirPath = usr.HomeDir + "/.dotfiles"

	if _, err := os.Stat(ConfigDirPath); err != nil {
		_ = os.Mkdir(ConfigDirPath, 0766)
	}
	if _, err := os.Stat(IanConfigPath); err != nil {
		_ = os.Mkdir(IanConfigPath, 0766)
		fmt.Printf("%s", []byte(`Welcome to Ian!
Ian is a simple tool to manage your development environment, repositories,
and projects.

Learn more about Ian at http://goian.io

To benefit from all of Ian’s feature’s, you’ll need to provide:
    - A working OS Package Manager (will set up if missing)
    - The full path of your repositories (example: /Users/thylong/repositories)
    - The path of your dotfiles Github repository (example: thylong/dotfiles)

`))
		if GetBoolUserInput("Do you want to set up Ian now? (Y/n)") == true {
			initVipers(false)
		} else {
			fmt.Println("You're ready to start using Ian. Note that if you try to use some of Ian's\nfeatures you'll be prompted for these details again.")
			return
		}
	}
	initVipers(false)
}

// initVipers return Vipers corresponding to Yaml config files.
// The soft argument determine if unexisting files should be written or not.
func initVipers(soft bool) {
	ConfigFilesPathes = make(map[string]string)
	Vipers = make(map[string]*viper.Viper)
	for _, ConfigFileName := range []string{"config", "env", "projects"} {
		ConfigFilesPathes[ConfigFileName] = IanConfigPath + fmt.Sprintf("%s.yml", ConfigFileName)
		Vipers[ConfigFileName] = initViper(ConfigFileName, soft)
	}
}

func initViper(viperName string, soft bool) (viperInstance *viper.Viper) {
	viperInstance = viper.New()
	viperInstance.SetConfigType("yaml")
	viperInstance.SetConfigName(viperName)
	viperInstance.AddConfigPath(IanConfigPath)

	if soft {
		return viperInstance
	}

	if _, err := os.Stat(fmt.Sprintf("%s/%s.yml", IanConfigPath, viperName)); err != nil {
		SetupConfigFile(viperName)
	}

	err := viperInstance.ReadInConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem with config file: %s Error: %s", viperName, err.Error())
		os.Exit(1)
	} else {
		configContent, _ := ioutil.ReadFile(ConfigFilesPathes[viperName])
		err = yaml.Unmarshal(configContent, &ConfigMap)
		if err != nil {
			fmt.Println("Unable to parse config file.")
			os.Exit(1)
		}
	}
	return viperInstance
}

// SetupConfigFile creates a config directory and the config file if not exists.
func SetupConfigFile(ConfigFileName string) {
	ConfigFilePath := ConfigFilesPathes[ConfigFileName]
	if _, err := os.Stat(ConfigFilePath); err != nil {
		configContent := GetConfigDefaultContent(ConfigFilePath)

		if ConfigFileName == "config" {
			repositoriesPathPrefix := "\nrepositories_path: "
			repositoriesPath := env.GenerateRepositoriesPath()
			configContent = append(configContent, fmt.Sprintf("%s%s", repositoriesPathPrefix, repositoriesPath)...)

			dotfilesRepositoryPrefix := "\ndotfiles_repository: "
			dotfilesRepository := env.GetDotfilesRepository()
			configContent = append(configContent, fmt.Sprintf("%s%s", dotfilesRepositoryPrefix, dotfilesRepository)...)
		}

		fmt.Printf("Creating %s\n", ConfigFileName)
		if err := ioutil.WriteFile(ConfigFilePath, configContent, 0766); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
			os.Exit(1)
		}
		return
	}
	fmt.Printf("Existing %s.yml found.\n", ConfigFileName)
}

// SetupConfigFiles creates a config directory and the config files if not exists.
func SetupConfigFiles() {
	for ConfigFileName := range ConfigFilesPathes {
		SetupConfigFile(ConfigFileName)
	}
}

// AppendToConfig takes a string as an argument
// and write it as new line(s) in the given conf file.
func AppendToConfig(lines string, confFilename string) {
	confPath := ConfigFilesPathes[confFilename]
	f, err := os.OpenFile(confPath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	if _, err = f.WriteString(lines); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// GetConfigDefaultContent returns the content of the default config.yml
// (As nothing is preset for now, this function actually returns an empty string)
func GetConfigDefaultContent(fileName string) []byte {
	return []byte{}
}

// GetPreset returns the content of the preset env.yml
func GetPreset(presetName string) []byte {
	return []byte{}
}

// UpdateYamlFile write a Viper content to a yaml file.
func UpdateYamlFile(fileFullPath string, fileContent map[string]interface{}) {
	out, err := yaml.Marshal(&fileContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update %s.\n", fileFullPath)
		os.Exit(1)
	}
	if err := ioutil.WriteFile(fileFullPath, out, 0766); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update %s.\n", fileFullPath)
		os.Exit(1)
	}
}

// GetUserInput ask question and return user input.
func GetUserInput(question string) string {
	fmt.Printf("%s: ", question)
	reader := bufio.NewReader(os.Stdin)
	if input, _ := reader.ReadString('\n'); input != "\n" && input != "" {
		return string(bytes.TrimSuffix([]byte(input), []byte("\n")))
	}
	return ""
}

// GetBoolUserInput ask question and return true if the user agreed otherwise false.
func GetBoolUserInput(question string) bool {
	in := GetUserInput(question)
	if strings.ToLower(in) == "y" && strings.ToLower(in) == "yes" && strings.ToLower(in) == "" {
		return true
	}
	return false
}

// SetupEnvFileWithPreset write an env file with the selected preset.
func SetupEnvFileWithPreset(preset string) {
	var Envcontent string
	switch preset {
	default:
		fmt.Fprint(os.Stderr, "Cannot find preset.")
		return
	case "1":
		Envcontent = GetSoftwareEngineerPreset()
	case "2":
		Envcontent = GetBackendDeveloperPreset()
	case "3":
		Envcontent = GetFrontendDeveloperPreset()
	case "4":
		Envcontent = GetOpsPreset()
	}

	confPath := ConfigFilesPathes["env"]
	f, err := os.OpenFile(confPath, os.O_CREATE|os.O_WRONLY, 0655)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	if _, err = f.WriteString(Envcontent); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// GetCustomCmds returns the commands defined in projects.yml
func GetCustomCmds(project string) (customCmds []*cobra.Command) {
	for confLineKey, confLineValue := range Vipers["projects"].GetStringMapString(project) {
		if strings.HasSuffix(confLineKey, "custom_cmd") {
			customCmdArgs := strings.Split(confLineValue, "=")
			customCmds = append(customCmds, &cobra.Command{
				Use:   strings.TrimSuffix(confLineKey, "_custom_cmd"),
				Short: customCmdArgs[0],
				Long:  customCmdArgs[0],
				Run: func(cmd *cobra.Command, args []string) {
					termCmd := exec.Command(customCmdArgs[1])
					command.ExecuteInteractiveCommand(termCmd)
				},
			})
		}
	}
	return customCmds
}

// GetProjects returns the projects defined in projects.yml as non-runnable []*cobra.cmd.
func GetProjects() (projectCmds map[string]*cobra.Command) {
	projectCmds = make(map[string]*cobra.Command)
	for _, project := range Vipers["projects"].AllKeys() {
		projectParams := Vipers["projects"].GetStringMapString(project)
		projectCmds[project] = &cobra.Command{
			Use:   project,
			Short: projectParams["description"],
			Long:  projectParams["description"],
		}
	}
	return projectCmds
}
