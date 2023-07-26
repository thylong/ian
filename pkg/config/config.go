// Copyright 2023 Théotime Levêque
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
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
	"path/filepath"
	"regexp"

	yaml "gopkg.in/yaml.v2"

	"github.com/spf13/viper"
	"github.com/thylong/ian/pkg/log"
)

// ConfigDirPath represents the path to config directory.
var ConfigDirPath string

// IanConfigPath represents the path to ian config directory.
var IanConfigPath string

// DotfilesDirPath represents the full path to the dotfiles directory.
var DotfilesDirPath string

// ConfigFilesPathes contains every config file pathes per filename.
var ConfigFilesPathes map[string]string

// Vipers contains all the initialized Vipers (config, env)
var Vipers map[string]*viper.Viper

// YamlConfigMap is used to marshal/unmarshal config file.
type YamlConfigMap struct {
	Managers     map[string]map[string]string `json:"managers"`
	Repositories map[string]string            `json:"repositories"`
	Setup        map[string][]string          `json:"setup"`
	Packages     map[string]map[string]string `json:"packages"`
}

// ConfigMap contains the config content.
var ConfigMap YamlConfigMap

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}

	ConfigDirPath = filepath.Join(usr.HomeDir, ".config")
	IanConfigPath = filepath.Join(ConfigDirPath, "ian")
	DotfilesDirPath = filepath.Join(usr.HomeDir, ".dotfiles")

	if _, err := os.Stat(ConfigDirPath); err != nil {
		_ = os.Mkdir(ConfigDirPath, 0766)
	}
	if _, err := os.Stat(IanConfigPath); err != nil {
		_ = os.Mkdir(IanConfigPath, 0766)
		log.Infoln(GetInitialSetupUsage())
	}

	ConfigFilesPathes = make(map[string]string)
	Vipers = make(map[string]*viper.Viper)
	initVipers()
}

// initVipers return Vipers corresponding to Yaml config files.
// The soft argument determine if unexisting files should be written or not.
func initVipers() {
	for _, ConfigFileName := range []string{"config", "env"} {
		configFilePath := filepath.Join(IanConfigPath, fmt.Sprintf("%s.yml", ConfigFileName))
		ConfigFilesPathes[ConfigFileName] = configFilePath
		Vipers[ConfigFileName] = initViper(ConfigFileName)
	}
}

// RefreshVipers is a helper called to refresh the configuration.
func RefreshVipers() {
	initVipers()
}

func initViper(viperName string) (viperInstance *viper.Viper) {
	viperInstance = viper.New()
	viperInstance.SetConfigType("yaml")
	viperInstance.SetConfigName(viperName)
	viperInstance.AddConfigPath(IanConfigPath)

	configFilePath := filepath.Join(IanConfigPath, fmt.Sprintf("%s.yml", viperName))
	if _, err := os.Stat(configFilePath); err != nil {
		SetupConfigFile(viperName)
	}

	err := viperInstance.ReadInConfig()
	if err != nil {
		log.Errorf("Problem with config file: %s: %s\n", viperName, err.Error())
		os.Exit(1)
	} else {
		configContent, _ := ioutil.ReadFile(ConfigFilesPathes[viperName])
		err = yaml.Unmarshal(configContent, &ConfigMap)
		if err != nil {
			log.Errorln("Unable to parse config file.")
			os.Exit(1)
		}
	}
	return viperInstance
}

// GetInitialSetupUsage returns the usage when using ian for the first time
func GetInitialSetupUsage() []byte {
	return []byte(`Welcome to Ian!
Ian is a simple tool to manage your development environment and repositories.

Learn more about Ian at http://goian.io

To benefit from all of Ian’s features, you’ll need to provide:
- The full path of your repositories (example: /Users/thylong/repositories)
- The path of your dotfiles Github repository (example: thylong/dotfiles)

`)
}

// SetupConfigFile creates a config directory and the config file if not exists.
func SetupConfigFile(ConfigFileName string) {
	ConfigFilePath := ConfigFilesPathes[ConfigFileName]
	if _, err := os.Stat(ConfigFilePath); err != nil {
		configContent := GetConfigDefaultContent(ConfigFilePath)

		if ConfigFileName == "config" {
			repositoriesPathPrefix := "repositories_path: "
			repositoriesPath := GenerateRepositoriesPath()
			configContent = append(configContent, fmt.Sprintf("%s%s", repositoriesPathPrefix, repositoriesPath)...)

			dotfilesRepositoryPrefix := "\ndotfiles:\n"
			dotfilesRepository := fmt.Sprintf("  repository: %s\n", GetDotfilesRepository())

			provider := "github"
			re := regexp.MustCompile("([a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9])\\.[a-zA-Z]{2,}")
			if len(re.FindStringSubmatch(dotfilesRepository)) > 1 {
				provider = re.FindStringSubmatch(dotfilesRepository)[1]
			}
			repositoryProvider := fmt.Sprintf("  provider: %s", provider)
			configContent = append(configContent, fmt.Sprintf("%s%s%s", dotfilesRepositoryPrefix, dotfilesRepository, repositoryProvider)...)
		}

		log.Infof("Creating %s\n", ConfigFileName)
		if err := ioutil.WriteFile(ConfigFilePath, configContent, 0766); err != nil {
			log.Errorf("%s\n", err)
			os.Exit(1)
		}
		return
	}
	log.Infof("Existing %s.yml found\n", ConfigFileName)
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
		log.Errorln("%s\n", err)
		os.Exit(1)
	}
	defer f.Close()

	if _, err = f.WriteString(lines); err != nil {
		log.Errorln("%s\n", err)
		os.Exit(1)
	}
}

// GetConfigDefaultContent returns the content of the default config.yml
// (As nothing is preset for now, this function actually returns an empty string)
func GetConfigDefaultContent(fileName string) []byte {
	return []byte{}
}

// UpdateYamlFile write a Viper content to a yaml file.
func UpdateYamlFile(fileFullPath string, fileContent map[string]interface{}) {
	out, err := yaml.Marshal(&fileContent)
	if err != nil {
		log.Errorf("Failed to update %s\n", fileFullPath)
		os.Exit(1)
	}
	if err := ioutil.WriteFile(fileFullPath, out, 0766); err != nil {
		log.Errorf("Failed to update %s\n", fileFullPath)
		os.Exit(1)
	}
}

// GenerateRepositoriesPath creates conf line containing the user's input.
func GenerateRepositoriesPath() string {
	reader := bufio.NewReader(os.Stdin)
	log.Infoln("\nEnter the full path to the parent directory of your repositories\n(leave blank to skip): ")
	if input, _ := reader.ReadString('\n'); input != "\n" && input != "" {
		return input
	}
	return ""
}

// GetDotfilesRepository creates conf line containing the user's input.
func GetDotfilesRepository() string {
	log.Infoln("\nEnter the full path to your dotfiles repository\n(leave blank to skip): ")
	reader := bufio.NewReader(os.Stdin)
	if input, _ := reader.ReadString('\n'); input != "\n" && input != "" {
		return string(bytes.TrimSuffix([]byte(input), []byte("\n")))
	}
	return ""
}

// GetDotfilesRepositoryPath returns the dotfiles repository path.
func GetDotfilesRepositoryPath() string {
	return Vipers["config"].GetStringMapString("dotfiles")["repository"]
}

// GetDefaultSaveMessage returns as a string the default save message.
func GetDefaultSaveMessage() string {
	return Vipers["config"].GetString("default_save_message")
}
