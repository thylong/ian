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
	"os/user"

	yaml "gopkg.in/yaml.v2"

	"github.com/spf13/viper"
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

	ConfigFilesPathes = make(map[string]string)
	Vipers = make(map[string]*viper.Viper)
	for _, ConfigFileName := range []string{"config", "env", "projects"} {
		ConfigFilesPathes[ConfigFileName] = IanConfigPath + fmt.Sprintf("%s.yml", ConfigFileName)
		Vipers[ConfigFileName] = initViper(ConfigFileName)
	}
}

func initViper(viperName string) (viperInstance *viper.Viper) {
	viperInstance = viper.New()
	viperInstance.SetConfigType("yaml")
	viperInstance.SetConfigName(viperName)
	viperInstance.AddConfigPath(IanConfigPath)

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
		viperInstance.WatchConfig()
	}
	return viperInstance
}

// SetupConfigFiles creates a config directory and the config files if not exists.
func SetupConfigFiles() {
	// Create .config dir if missing.
	if _, err := os.Stat(ConfigDirPath); err != nil {
		_ = os.Mkdir(ConfigDirPath, 0766)
	}
	// Create .config/ian dir if missing.
	if _, err := os.Stat(IanConfigPath); err != nil {
		_ = os.Mkdir(IanConfigPath, 0766)
	}

	for ConfigFileName, ConfigFilePath := range ConfigFilesPathes {
		if _, err := os.Stat(ConfigFilePath); err != nil {
			configContent := GetConfigDefaultContent(ConfigFilePath)

			repositoriesPath := generateRepositoriesPath()
			configContent = append(configContent, repositoriesPath...)

			GithubUsername := generateGithubUsername()
			configContent = append(configContent, GithubUsername...)

			fmt.Printf("Creating %s", ConfigFileName)
			if err := ioutil.WriteFile(ConfigFilePath, configContent, 0766); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
				os.Exit(1)
			}
			return
		}
	}
	fmt.Print("Config files found.")
}

// GetConfigDefaultContent returns the content of the default config.yml
func GetConfigDefaultContent(fileName string) []byte {
	return []byte{}
}

func generateRepositoriesPath() (repositoriesPathConf string) {
	repositoriesPathConf = "\nrepositories_path: "
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Ian allows you manage all your Github repositories")
	fmt.Print("Insert up the fullpath to the parent directory of all your repositories, otherwise leave blank: ")
	if fullPathToRepositories, _ := reader.ReadString('\n'); fullPathToRepositories != "\n" {
		return repositoriesPathConf + fullPathToRepositories
	}
	return repositoriesPathConf
}

func generateGithubUsername() (GithubUsernameConf string) {
	GithubUsernameConf = "\ngithub_username: "
	fmt.Println("Insert your Github nickname: ")
	reader := bufio.NewReader(os.Stdin)
	if nickname, _ := reader.ReadString('\n'); nickname != "\n" {
		Vipers["config"].Set("github_username", string(bytes.TrimSuffix([]byte(nickname), []byte("\n"))))
		return GithubUsernameConf + string(bytes.TrimSuffix([]byte(nickname), []byte("\n")))
	}
	return GithubUsernameConf
}
