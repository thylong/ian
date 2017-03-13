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
	"regexp"

	"github.com/thylong/ian/backend/command"
	pm "github.com/thylong/ian/backend/package-managers"
)

// GetInfos returns env infos
func GetInfos() {
	IPCheckerURL := "http://httpbin.org/ip"

	resp, err := http.Get(IPCheckerURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
	}
	content, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var jsonContent map[string]string
	err = json.Unmarshal(content, &jsonContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
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
			fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
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
			fmt.Fprintf(os.Stderr, "couldn't move %s !", src)
			os.Exit(1)
		}
		if err := os.Symlink(dst, src); err != nil {
			fmt.Fprintf(os.Stderr, "couldn't symlink %s !", src)
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
		fmt.Fprintf(os.Stderr, "%s repository doesn't exists or is not reachable.", repositoryURL)
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
		fmt.Fprint(os.Stderr, "Cannot interact with Git")
	}

	commitCmd := exec.Command("git", "commit", "-m", message)
	commitCmd.Dir = dotfilesDirPath
	if err = command.MustExecuteCommand(commitCmd); err != nil {
		fmt.Fprint(os.Stderr, "Cannot create a commit")
	}

	termCmd := exec.Command("git", "push", "--force", "origin", "master")
	termCmd.Dir = dotfilesDirPath
	if err = command.MustExecuteCommand(termCmd); err != nil {
		fmt.Fprint(os.Stderr, "Cannot push to repository")
	}
}

// SetupPackages installs listed CLI packages.
func SetupPackages(PackageManager pm.PackageManager, packages []string) {
	fmt.Println("Installing packages...")

	if len(packages) == 0 {
		fmt.Println("No packages to install")
		return
	}

	for _, packageToInstall := range packages {
		PackageManager.Install(packageToInstall)
	}
}

// SetupDotFiles ask and retrieve a dotfiles repository.
func SetupDotFiles(dotfilesRepository string, dotfilesDirPath string) {
	usr, _ := user.Current()
	if _, err := os.Stat(usr.HomeDir + "/.dotfiles"); err != nil && dotfilesRepository != "" {
		termCmd := exec.Command("git", "clone", "-v", "https://github.com/"+dotfilesRepository+".git", dotfilesDirPath)
		termCmd.Stdout = os.Stdout
		termCmd.Stdin = os.Stdin
		termCmd.Stderr = os.Stderr
		termCmd.Run()

		re := regexp.MustCompile(".git$")

		files, _ := ioutil.ReadDir(usr.HomeDir + "/.dotfiles")
		for _, f := range files {
			if re.MatchString(f.Name()) {
				continue
			}

			if _, err := os.Stat(usr.HomeDir + "/" + f.Name()); err != nil {
				err := os.Symlink(usr.HomeDir+"/.dotfiles/"+f.Name(), usr.HomeDir+"/"+f.Name())
				if err != nil {
					fmt.Fprint(os.Stderr, err)
				}
			}
		}
	} else {
		fmt.Println("Skipping dotfiles configuration.")
	}
}

// GenerateRepositoriesPath creates conf line containing the user's input.
func GenerateRepositoriesPath() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Ian allows you manage all your Github local repositories")
	fmt.Print("Insert the full path to the parent directory of your repositories, otherwise leave blank: ")
	if fullPathToRepositories, _ := reader.ReadString('\n'); fullPathToRepositories != "\n" {
		return fullPathToRepositories
	}
	return ""
}

// GetDotfilesRepository creates conf line containing the user's input.
func GetDotfilesRepository() string {
	fmt.Println("Path to your dotfiles repository: ")
	reader := bufio.NewReader(os.Stdin)
	if input, _ := reader.ReadString('\n'); input != "\n" {
		// Vipers["config"].Set("dotfiles_repository", string(bytes.TrimSuffix([]byte(input), []byte("\n"))))
		return string(bytes.TrimSuffix([]byte(input), []byte("\n")))
	}
	return ""
}
