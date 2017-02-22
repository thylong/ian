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
	"regexp"

	yaml "gopkg.in/yaml.v2"

	"github.com/spf13/viper"
	"github.com/thylong/ian/backend/command"
)

// ConfigPath represents the path to config directory.
var ConfigPath string

// IanConfigPath represents the path to ian config directory.
var IanConfigPath string

// ConfigFullPath represents the path to config file.
var ConfigFullPath string

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
	// Init and keep track of config.
	usr, err := user.Current()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	ConfigPath = usr.HomeDir + "/.config"
	IanConfigPath = ConfigPath + "/ian/"
	ConfigFullPath = IanConfigPath + "config.yml"
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(IanConfigPath)

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem with config file: %s", err.Error())
		os.Exit(1)
	} else {
		configContent, _ := ioutil.ReadFile(ConfigFullPath)
		err = yaml.Unmarshal(configContent, &ConfigMap)
		if err != nil {
			fmt.Println("Unable to parse config file.")
			return
		}
		viper.WatchConfig()
	}
}

// SetupConfigFile creates a config directory and store the default config file
// If not exists.
func SetupConfigFile() {
	// Create .config dir if missing.
	_, err := os.Stat(ConfigPath)
	if err != nil {
		err = os.Mkdir(ConfigPath, 0766)
	}
	// Create .config/ian dir if missing.
	_, err = os.Stat(IanConfigPath)
	if err != nil {
		err = os.Mkdir(IanConfigPath, 0766)
	}

	// Create config.yml file
	_, err = os.Stat(ConfigFullPath)
	if err != nil {
		configContent := GetConfigDefaultContent()

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Ian allows you manage all your Github repositories")
		fmt.Print("leave blank if you don't want to benefit of this feature.")
		fmt.Print("Otherwise, set up the fullpath to the parent directory of all your repositories: ")
		if fullPathToRepositories, _ := reader.ReadString('\n'); fullPathToRepositories != "\n" {
			repositoriesPathConfigLine := "\nrepositories_path: " + fullPathToRepositories
			configContent = append(configContent, repositoriesPathConfigLine...)
		}

		fmt.Printf("Creating %s", ConfigFullPath)
		err := ioutil.WriteFile(ConfigFullPath, configContent, 0766)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
			os.Exit(1)
		}
		return
	}
	fmt.Print("Config file found.")
}

// SetupDotFiles ask for a Github nickname and retrieve the dotfiles repo
// (the repository has to be public).
func SetupDotFiles() {
	reader := bufio.NewReader(os.Stdin)
	usr, err := user.Current()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	dotfilesPath := usr.HomeDir + "/.dotfiles"

	fmt.Print("Ian allows you to automatically import your dotfiles from a Github repository,")
	fmt.Print("leave blank if you don't want to benefit of this feature.")
	fmt.Print("Github nickname: ")
	if nickname, _ := reader.ReadString('\n'); nickname != "\n" {
		nickname = string(bytes.TrimSuffix([]byte(nickname), []byte("\n")))
		termCmd := exec.Command("git", "clone", "-v", "https://github.com/"+nickname+"/dotfiles", dotfilesPath)
		command.ExecuteCommand(termCmd)

		re := regexp.MustCompile(".git$")

		usr, _ := user.Current()
		files, _ := ioutil.ReadDir(usr.HomeDir + "/.dotfiles")
		for _, f := range files {
			if re.MatchString(f.Name()) {
				continue
			}
			if err := os.Symlink(usr.HomeDir+"/.dotfiles/"+f.Name(), usr.HomeDir+"/"+f.Name()); err != nil {
				fmt.Fprint(os.Stderr, err)
			}
		}
	} else {
		fmt.Print("Skipping dotfiles configuration.")
	}
}

// GetConfigDefaultContent returns the content of the default config.yml
func GetConfigDefaultContent() []byte {
	return []byte(`managers:
  pip:
    base_url: https://pypi.org/project/
    install_cmd: install
    uninstall_cmd: uninstall
  npm:
    base_url: https://www.npmjs.com/package/
    install_cmd: install -g
    uninstall_cmd: uninstall -g
  gem:
    base_url: https://rubygems.org/gems/
    install_cmd: install
    uninstall_cmd: uninstall
  composer: # alias packagist
    base_url: https://packagist.org/packages/
    install_cmd: global install
    uninstall_cmd: global remove
repositories: {}
projects:
  cabu:
    repository: thylong/cabu
  ian:
    repository: thylong/ian
  regexrace:
    db_cmd: mongo localhost
    deploy_cmd: bash deploy.sh
    health: /status
    repository: thylong/regexrace
    rollback_cmd: bash rollback.sh
    url: http://regexrace.org
  thylong:
    health: /
    repository: thylong/thylong.github.io
    url: http://thylong.com
setup:
  cli_packages:
  - httpie
  - fish
  - keybase
  - mongodb
  - lynx
  - node
  - nmap
  - python
  - python3
  - rsyslog
  - cmake
  - ruby
  - tree
  - cask
  - tmux
  - wget
  - reattach-to-user-namespace
  gui_packages:
  - appcleaner
  - atom
  - caffeine
  - charles
  - dash
  - dashlane
  - docker
  - filezilla
  - firefox
  - google-chrome
  - iterm2
  - jadengeller-helium
  - keka
  - libreoffice
  - mediainfo
  - robomongo
  - skype
  - slack
  - spectacle
  - spotify
  - steam
  - torbrowser
  - tunnelblick
  - utorrent
  - vagrant
  - virtualbox
  - vlc
  - wireshark
  requirements:
  - git
  - gcloud
packages:
  baily-cli:
    cmd: baily-cli
    description: baily-cli
    type: npm
`)
}
