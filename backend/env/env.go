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
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/afero"
	"github.com/thylong/ian/backend/command"
)

// AppFs is a wrapper to OS package
var AppFs = afero.NewOsFs()

var httpGet = http.Get

var execCommand = exec.Command

// IPCheckerURL is the endpoint to call to get IP data
var IPCheckerURL = "http://httpbin.org/ip"

// ErrJSONPayloadInvalidFormat is returned when the JSON payload format is invalid
var ErrJSONPayloadInvalidFormat = fmt.Errorf("%v %s", color.RedString("Error:"), errors.New("Invalid JSON format"))

// ErrOperationNotPermitted is returned when trying create or write without permissions
var ErrOperationNotPermitted = fmt.Errorf("%v %s", color.RedString("Error:"), errors.New("Operation not permitted"))

// ErrCannotMoveDotfile is returned when trying create or write without permissions
var ErrCannotMoveDotfile = fmt.Errorf("%v couldn't move dotfile", color.RedString("Error:"))

// ErrCannotSymlink is returned when trying to create a Symlink and fails
var ErrCannotSymlink = fmt.Errorf("%v couldn't create symlink", color.RedString("Error:"))

// ErrCannotInteractWithGit is returned when trying to interact with Git
var ErrCannotInteractWithGit = fmt.Errorf("%v Cannot interact with Git\n", color.RedString("Error:"))

// ErrHTTPError is returned when failing to reach an endpoint with HTTP
var ErrHTTPError = fmt.Errorf("%v Cannot reach endpoint", color.RedString("Error:"))

// ErrDotfilesRepository is returned when failing to stat a repository
var ErrDotfilesRepository = fmt.Errorf("%v dotfiles repository doesn't exists or is not reachable", color.RedString("Error:"))

// Add a package in env.yml.
func Add(packageManager string, packages []string) (NewPMList []string, err error) {
	return NewPMList, err
}

// Remove a package in env.yml.
func Remove(packageManager string, packages []string) (NewPMList []string, err error) {
	return NewPMList, err
}

// Describe returns env description.
func Describe() (err error) {
	resp, err := httpGet(IPCheckerURL)
	if err != nil {
		return ErrHTTPError
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%v %s", color.RedString("Error:"), err)
	}
	defer resp.Body.Close()

	var jsonContent map[string]string
	err = json.Unmarshal(content, &jsonContent)
	if err != nil {
		return ErrJSONPayloadInvalidFormat
	}

	command.ExecuteCommand(execCommand("hostinfo"))
	fmt.Printf("\nExternal IP: %s\n\n", jsonContent["origin"])
	fmt.Print("Uptime: ")
	command.ExecuteCommand(execCommand("uptime"))

	return nil
}

// Save persists the dotfiles in distant repository.
func Save(dotfilesDirPath string, dotfilesRepository string, defaultSaveMessage string, dotfilesToSave []string) (err error) {
	if err = EnsureDotfilesDir(dotfilesDirPath); err != nil {
		return err
	}
	if err = ImportIntoDotfilesDir(dotfilesToSave, dotfilesDirPath); err != nil {
		return err
	}
	if err = EnsureDotfilesRepository(dotfilesRepository, dotfilesDirPath); err != nil {
		return err
	}
	if err = PersistDotfiles(defaultSaveMessage, dotfilesDirPath); err != nil {
		return err
	}
	return nil
}

// EnsureDotfilesDir create the ~/.dotfiles directory if not exists.
func EnsureDotfilesDir(dotfilesDirPath string) (err error) {
	dotfilesDirPath = filepath.Dir(dotfilesDirPath)
	if _, err := AppFs.Stat(dotfilesDirPath); err != nil {
		err = AppFs.Mkdir(dotfilesDirPath, 0766)
		if err != nil {
			return ErrOperationNotPermitted
		}
		command.ExecuteCommand(execCommand("git", "init"))
		GitIgnorePath := filepath.Join(dotfilesDirPath, ".gitignore")
		ioutil.WriteFile(GitIgnorePath, []byte(".ssh\n.netrc"), 0766)
	}
	return nil
}

// ImportIntoDotfilesDir moves dotfiles into dotfiles directory and create symlinks.
func ImportIntoDotfilesDir(dotfilesToSave []string, dotfilesDirPath string) (err error) {
	usr, _ := user.Current()

	if len(dotfilesToSave) == 0 {
		files, _ := ioutil.ReadDir(usr.HomeDir)
		for _, file := range files {
			if strings.HasPrefix(file.Name(), ".") && file.Name() != ".ssh" && file.Name() != ".bash_history" && file.Name() != ".Trash" {
				dotfilesToSave = append(dotfilesToSave, file.Name())
			}
		}
	}
	for _, dotfileToSave := range dotfilesToSave {
		src := filepath.Join(usr.HomeDir, dotfileToSave)
		dst := filepath.Join(dotfilesDirPath, dotfileToSave)

		if err := MoveFile(src, dst); err != nil {
			return ErrCannotMoveDotfile
		}
		if err := os.Symlink(dst, src); err != nil {
			return ErrCannotSymlink
		}
	}
	fmt.Printf("Moved dotfiles in %s directory.\n", dotfilesDirPath)
	return nil
}

// EnsureDotfilesRepository create Dotfiles repository if not exists.
func EnsureDotfilesRepository(dotfilesRepository string, dotfilesDirPath string) (err error) {
	if dotfilesRepository == "" {
		dotfilesRepository = GetDotfilesRepository()
	}

	repositoryURL := fmt.Sprintf("git@github.com:%s.git", dotfilesRepository)
	lsRemoteCmd := execCommand("git", "ls-remote", repositoryURL)
	lsRemoteCmd.Dir = dotfilesDirPath

	if err := command.MustExecuteCommand(lsRemoteCmd); err != nil {
		fmt.Println(err)
		return ErrDotfilesRepository
	}
	return nil
}

// PersistDotfiles local dotfiles to remote.
func PersistDotfiles(message string, dotfilesDirPath string) (err error) {
	if len(message) == 0 {
		message = "Update dotfiles"
	}

	addCmd := execCommand("git", "add", "-A")
	addCmd.Dir = dotfilesDirPath
	if err = command.MustExecuteCommand(addCmd); err != nil {
		return ErrCannotInteractWithGit
	}

	commitCmd := execCommand("git", "commit", "-m", message)
	commitCmd.Dir = dotfilesDirPath
	if err = command.MustExecuteCommand(commitCmd); err != nil {
		return fmt.Errorf("%v Cannot create a commit\n", color.RedString("Error:"))
	}

	pushCmd := execCommand("git", "push", "--force", "origin", "master")
	pushCmd.Dir = dotfilesDirPath
	if err = command.MustExecuteCommand(pushCmd); err != nil {
		return fmt.Errorf("%v Cannot push to repository.\n", color.RedString("Error:"))
	}
	return nil
}

// GenerateRepositoriesPath creates conf line containing the user's input.
func GenerateRepositoriesPath() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nEnter the full path to the parent directory of your repositories\n(leave blank to skip): ")
	if input, _ := reader.ReadString('\n'); input != "\n" && input != "" {
		return input
	}
	return ""
}

// GetDotfilesRepository creates conf line containing the user's input.
func GetDotfilesRepository() string {
	fmt.Print("\nEnter the full path to your dotfiles repository\n(leave blank to skip): ")
	reader := bufio.NewReader(os.Stdin)
	if input, _ := reader.ReadString('\n'); input != "\n" && input != "" {
		return string(bytes.TrimSuffix([]byte(input), []byte("\n")))
	}
	return ""
}

// GetInitialSetupUsage returns the usage when using ian for the first time
func GetInitialSetupUsage() []byte {
	return []byte(`Welcome to Ian!
Ian is a simple tool to manage your development environment, repositories,
and projects.

Learn more about Ian at http://goian.io

To benefit from all of Ian’s features, you’ll need to provide:
- The full path of your repositories (example: /Users/thylong/repositories)
- The path of your dotfiles Github repository (example: thylong/dotfiles)

`)
}
