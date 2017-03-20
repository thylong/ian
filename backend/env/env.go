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
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"

	"github.com/fatih/color"
	"github.com/thylong/ian/backend/command"
)

// Describe returns env description.
func Describe() {
	IPCheckerURL := "http://httpbin.org/ip"

	resp, err := http.Get(IPCheckerURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v %s.", color.RedString("Error:"), err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var jsonContent map[string]string
	err = json.Unmarshal(content, &jsonContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v %s.", color.RedString("Error:"), err)
		return
	}

	command.ExecuteCommand(exec.Command("hostinfo"))
	fmt.Printf("\nExternal IP: %s\n\n", jsonContent["origin"])
	fmt.Print("Uptime: ")
	command.ExecuteCommand(exec.Command("uptime"))
}

// Save persists the dotfiles in distant repository.
func Save(dotfilesDirPath string, dotfilesRepository string, defaultSaveMessage string, dotfilesToSave []string) {
	EnsureDotfilesDir(dotfilesDirPath)
	ImportIntoDotfilesDir(dotfilesToSave, dotfilesDirPath)
	EnsureDotfilesRepository(dotfilesRepository, dotfilesDirPath)
	PushDotfiles(defaultSaveMessage, dotfilesDirPath)
}

// EnsureDotfilesDir create the ~/.dotfiles directory if not exists.
func EnsureDotfilesDir(dotfilesDirPath string) {
	if _, err := os.Stat(dotfilesDirPath); err != nil {
		err = os.Mkdir(dotfilesDirPath, 0766)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v %s.", color.RedString("Error:"), err)
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

		if err := MoveFile(src, dst); err != nil {
			fmt.Fprintf(os.Stderr, "%v couldn't move %s !", color.RedString("Error:"), src)
			os.Exit(1)
		}
		if err := os.Symlink(dst, src); err != nil {
			fmt.Fprintf(os.Stderr, "%v couldn't symlink %s !", color.RedString("Error:"), src)
			os.Exit(1)
		}
	}
	fmt.Printf("Moved dotfiles in %s directory.\n", dotfilesDirPath)
}

// EnsureDotfilesRepository create Dotfiles repository if not exists.
func EnsureDotfilesRepository(dotfilesRepository string, dotfilesDirPath string) {
	if dotfilesRepository == "" {
		dotfilesRepository = GetDotfilesRepository()
	}
	repositoryURL := fmt.Sprintf("git@github.com:%s.git", dotfilesRepository)
	termCmd := exec.Command("git", "ls-remote", repositoryURL)
	termCmd.Dir = dotfilesDirPath

	if err := command.MustExecuteCommand(termCmd); err != nil {
		fmt.Fprintf(os.Stderr, "%v %s repository doesn't exists or is not reachable.", color.RedString("Error:"), repositoryURL)
		os.Exit(1)
	}
}

// PushDotfiles local dotfiles to remote.
func PushDotfiles(message string, dotfilesDirPath string) {
	var err error
	if message == "" {
		message = "Update dotfiles"
	}

	addCmd := exec.Command("git", "add", "-A")
	addCmd.Dir = dotfilesDirPath
	if err = command.MustExecuteCommand(addCmd); err != nil {
		fmt.Fprintf(os.Stderr, "%v Cannot interact with Git.\n", color.RedString("Error:"))
		return
	}

	commitCmd := exec.Command("git", "commit", "-m", message)
	commitCmd.Dir = dotfilesDirPath
	if err = command.MustExecuteCommand(commitCmd); err != nil {
		fmt.Fprintf(os.Stderr, "%v Cannot create a commit.\n", color.RedString("Error:"))
		return
	}

	termCmd := exec.Command("git", "push", "--force", "origin", "master")
	termCmd.Dir = dotfilesDirPath
	if err = command.MustExecuteCommand(termCmd); err != nil {
		fmt.Fprintf(os.Stderr, "%v Cannot push to repository.\n", color.RedString("Error:"))
	}
}

// GenerateRepositoriesPath creates conf line containing the user's input.
func GenerateRepositoriesPath() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the full path to the parent directory of your repositories\n(leave blank to skip): ")
	if input, _ := reader.ReadString('\n'); input != "\n" && input != "" {
		return input
	}
	return ""
}

// GetDotfilesRepository creates conf line containing the user's input.
func GetDotfilesRepository() string {
	fmt.Print("Enter the full path to your dotfiles Github repository\n(leave blank to skip): ")
	reader := bufio.NewReader(os.Stdin)
	if input, _ := reader.ReadString('\n'); input != "\n" && input != "" {
		return string(bytes.TrimSuffix([]byte(input), []byte("\n")))
	}
	return ""
}
