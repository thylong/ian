// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup ian working environment",
	Long: `Ian requires you to be able to interact with Github through Git CLI.

    With projects subcommand being one of the core function of Ian, setup will
    install what is necessessary to deploy on GCE.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting setup")
		fmt.Println("====================")

		if _, err := os.Stat("/usr/local/bin/brew"); err != nil {
			log.Fatal("Missing homebrew !")
		}

		setupConfigFile()
		setupDotFiles()
		setupBrewPackages()
		setupCask()
		setupCaskPackages()

		fmt.Println("====================")
		fmt.Println("Ending setup.")
	},
}

func init() {
	RootCmd.AddCommand(setupCmd)
}

func setupConfigFile() {
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
		fmt.Println("Creating ", ConfigFullPath)
		err := ioutil.WriteFile(ConfigFullPath, getConfigDefaultContent(), 0766)
		if err != nil {
			panic(err.Error())
		}
		return
	}
	fmt.Println("Config file found.")
}

func setupBrewPackages() {
	fmt.Println("Installing brew packages...")
	brewPackages := viper.GetStringMap("setup")["brew_packages"]

	if brewPackages == nil {
		fmt.Println("No brew packages to install")
		return
	}

	for _, brewPackage := range brewPackages.([]interface{}) {
		executeCommand(exec.Command("/usr/local/bin/brew", "install", brewPackage.(string)))
	}
}

// This function has often no output.
func setupCask() {
	fmt.Println("Installing cask...")
	if _, err := os.Stat("/usr/local/bin/cask"); err != nil {
		executeCommand(exec.Command("brew", "tap", "caskroom/cask"))
	} else {
		fmt.Println("cask already installed, skipping...")
	}
}

func setupCaskPackages() {
	fmt.Println("Installing cask packages...")
	caskPackages := viper.GetStringMap("setup")["cask_packages"]

	if caskPackages == nil {
		fmt.Println("No cask packages to install")
		return
	}

	for _, caskPackage := range caskPackages.([]interface{}) {
		executeCommand(exec.Command("/usr/local/bin/brew", "cask", "install", caskPackage.(string)))
	}
}

// setupDotFiles ask for a Github nickname and retrieve the dotfiles repo
// (the repository has to be public).
func setupDotFiles() {
	reader := bufio.NewReader(os.Stdin)
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dotfilesPath := usr.HomeDir + "/.dotfiles"

	fmt.Print("Ian allows you to automatically import your dotfiles from a Github repository,")
	fmt.Print("leave blank if you don't want to benefit of this feature.")
	fmt.Print("Github nickname: ")
	if nickname, _ := reader.ReadString('\n'); nickname != "\n" {
		nickname = string(bytes.TrimSuffix([]byte(nickname), []byte("\n")))
		termCmd := exec.Command("git", "clone", "-v", "https://github.com/"+nickname+"/dotfiles", dotfilesPath)
		executeCommand(termCmd)

		re := regexp.MustCompile(".git$")

		usr, _ := user.Current()
		files, _ := ioutil.ReadDir(usr.HomeDir + "/.dotfiles")
		for _, f := range files {
			if re.MatchString(f.Name()) {
				continue
			}
			if err := os.Symlink(usr.HomeDir+"/.dotfiles/"+f.Name(), usr.HomeDir+"/"+f.Name()); err != nil {
				log.Debug(err)
			}
		}
	} else {
		fmt.Println("Skipping dotfiles configuration.")
	}
}

// getConfigDefaultContent returns the content of the default config.yml
func getConfigDefaultContent() []byte {
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
  brew_packages:
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
  cask_packages:
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
