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

package env

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"github.com/thylong/ian/pkg/command"
	"github.com/thylong/ian/pkg/config"
	"github.com/thylong/ian/pkg/log"
)

// AppFs is a wrapper to OS package
var AppFs = afero.NewOsFs()

var httpGet = http.Get

var execCommand = exec.Command

// IPCheckerURL is the endpoint to call to get IP data
var IPCheckerURL = "http://httpbin.org/ip"

// Add a package in env.yml.
func Add(packageManager string, packages []string) (NewPMList []string, err error) {
	return NewPMList, err
}

// Remove a package in env.yml.
func Remove(packageManager string, packages []string) (NewPMList []string, err error) {
	return NewPMList, err
}

// Save persists the dotfiles in distant repository.
func Save(dotfilesToSave []string) (err error) {
	if err = EnsureDotfilesDir(config.DotfilesDirPath); err != nil {
		return err
	}
	if err = ImportIntoDotfilesDir(dotfilesToSave, config.DotfilesDirPath); err != nil {
		return err
	}
	if err = EnsureDotfilesRepository(config.GetDotfilesRepositoryPath(), config.DotfilesDirPath); err != nil {
		return err
	}
	if err = PersistDotfiles(config.GetDefaultSaveMessage(), config.DotfilesDirPath); err != nil {
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
	log.Infof("Moved dotfiles in %s directory\n", dotfilesDirPath)
	return nil
}

// EnsureDotfilesRepository create Dotfiles repository if not exists.
func EnsureDotfilesRepository(dotfilesRepository string, dotfilesDirPath string) (err error) {
	if dotfilesRepository == "" {
		dotfilesRepository = config.GetDotfilesRepository()
	}

	repositoryURL := fmt.Sprintf("git@github.com:%s.git", dotfilesRepository)
	lsRemoteCmd := execCommand("git", "ls-remote", repositoryURL)
	lsRemoteCmd.Dir = dotfilesDirPath

	if err := command.MustExecuteCommand(lsRemoteCmd); err != nil {
		log.Errorln(err)
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
		return errors.New("Cannot create a commit")
	}

	pushCmd := execCommand("git", "push", "--force", "origin", "master")
	pushCmd.Dir = dotfilesDirPath
	if err = command.MustExecuteCommand(pushCmd); err != nil {
		return errors.New("Cannot push to repository")
	}
	return nil
}

// AddPackagesToEnvFile adds packages to the env.yml file.
func AddPackagesToEnvFile(packageManagerName string, packages []string) {
	envContent := config.Vipers["env"].AllSettings()
	pmContent := config.Vipers["env"].GetStringSlice(packageManagerName)
	contains := func(e []string, c string) bool {
		for _, s := range e {
			if s == c {
				return true
			}
		}
		return false
	}
	for _, p := range packages {
		if !contains(pmContent, p) {
			pmContent = append(pmContent, p)
		}
	}

	envContent[packageManagerName] = pmContent
	config.UpdateYamlFile(
		config.ConfigFilesPathes["env"],
		envContent,
	)
}

// RemovePackagesFromEnvFile removes packages from the env.yml file.
func RemovePackagesFromEnvFile(packageManagerName string, packages []string) {
	envContent := config.Vipers["env"].AllSettings()
	pmContent := config.Vipers["env"].GetStringSlice(packageManagerName)
	contains := func(e []string, c string) bool {
		for _, s := range e {
			if s == c {
				return true
			}
		}
		return false
	}
	for i, p := range packages {
		if contains(pmContent, p) {
			pmContent = append(pmContent[:i], pmContent[i+1:]...)
		}
	}

	envContent[packageManagerName] = pmContent
	config.UpdateYamlFile(
		config.ConfigFilesPathes["env"],
		envContent,
	)
}
