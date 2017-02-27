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

package env

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"

	"github.com/thylong/ian/backend/command"
	"github.com/thylong/ian/backend/lib"
)

// GetInfos returns env infos
func GetInfos() {
	IPCheckerURL := "http://httpbin.org/ip"

	resp, err := http.Get(IPCheckerURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error : %s", err.Error())
	}
	content, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var jsonContent map[string]string
	err = json.Unmarshal(content, &jsonContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error : %s", err.Error())
		return
	}

	command.ExecuteCommand(exec.Command("hostinfo"))
	fmt.Println("external_ip :", jsonContent["origin"])
	fmt.Print("uptime :")
	command.ExecuteCommand(exec.Command("uptime"))
}

// EnsureDotfilesDir create the ~/.dotfiles directory if not exists.
func EnsureDotfilesDir(dotfilesDirPath string) {
	if _, err := os.Stat(dotfilesDirPath); err != nil {
		err = os.Mkdir(dotfilesDirPath, 0766)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error : %s", err.Error())
		}

		command.ExecuteCommand(exec.Command("git", "init"))
		GitIgnorePath := fmt.Sprintf("%s/.gitignore", dotfilesDirPath)
		_ = ioutil.WriteFile(GitIgnorePath, []byte(".ssh\n.netrc"), 0766)
		return
	}
}

// ImportIntoDotfilesDir moves dotfiles into dotfiles directory and create symlinks.
func ImportIntoDotfilesDir(dotfilesToSave []string, dotfilesDirPath string) {
	usr, _ := user.Current()
	for _, dotfileToSave := range dotfilesToSave {
		src := fmt.Sprintf("%s/%s", usr.HomeDir, dotfileToSave)
		dst := fmt.Sprintf("%s/%s", dotfilesDirPath, dotfileToSave)

		if err := lib.MoveFile(src, dst); err != nil {
			panic(fmt.Sprintf("couldn't move %s !", src))
		}
		if err := os.Symlink(dst, src); err != nil {
			panic(fmt.Sprintf("couldn't symlink %s !", src))
		}
	}
	fmt.Printf("Moved dotfiles in %s directory.\n", dotfilesDirPath)
}

// EnsureDotfilesRepository create Dotfiles repository if not exists.
func EnsureDotfilesRepository(githubUser string, dotfilesDirPath string) {
	repositoryURL := fmt.Sprintf("git@github.com:%s/testong.git", githubUser)
	termCmd := exec.Command("git", "ls-remote", repositoryURL)
	termCmd.Dir = dotfilesDirPath

	if err := command.MustExecuteCommand(termCmd); err != nil {
		fmt.Printf("%s repository doesn't exists or is not reachable.", repositoryURL)
		os.Exit(1)
	}
}

// PushDotfiles local dotfiles to remote.
func PushDotfiles(message string, dotfilesDirPath string) {
	if message == "" {
		message = "Update dotfiles"
	}

	addCmd := exec.Command("git", "add", "-A")
	addCmd.Dir = dotfilesDirPath
	command.ExecuteCommand(addCmd)

	commitCmd := exec.Command("git", "commit", "-m", message)
	commitCmd.Dir = dotfilesDirPath
	command.ExecuteCommand(commitCmd)

	termCmd := exec.Command("git", "push", "origin", "master")
	termCmd.Dir = dotfilesDirPath
	command.ExecuteCommand(termCmd)
}
