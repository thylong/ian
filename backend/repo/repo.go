package repo

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/thylong/ian/backend/command"
	"github.com/thylong/ian/backend/config"
	"github.com/thylong/ian/backend/log"
)

var execCommand = exec.Command

// List local repositories
func List() error {
	termCmd := execCommand("ls")
	log.Infof("repositories_path: %s\n", config.Vipers["config"].GetString("repositories_path"))
	termCmd.Dir = config.Vipers["config"].GetString("repositories_path")

	return command.ExecuteCommand(termCmd)
}

// Clone local repository
func Clone(repository string) error {
	termCmd := execCommand("git", "clone", "-v", repository)
	termCmd.Dir = config.Vipers["config"].GetString("repositories_path")

	command.ExecuteInteractiveCommand(termCmd)
	return nil
}

// Clean given repository
func Clean(repository string) error {
	termCmd := execCommand("git", "clean", "-dffx", repository)
	termCmd.Dir = config.Vipers["config"].GetString("repositories_path")

	return command.ExecuteCommand(termCmd)
}

// UpdateAll local repositories
func UpdateAll() {
	files, err := ioutil.ReadDir(config.Vipers["config"].GetString("repositories_path"))
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
	for _, file := range files {
		if file.IsDir() {
			UpdateOne(file.Name())
		}
	}
}

// UpdateOne local repository
func UpdateOne(repository string) error {
	termCmd := execCommand("git", "fetch", repository)
	termCmd.Dir = config.Vipers["config"].GetString("repositories_path")

	return command.ExecuteCommand(termCmd)
}

// UpgradeAll local repositories
func UpgradeAll() {
	files, err := ioutil.ReadDir(config.Vipers["config"].GetString("repositories_path"))
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
	for _, file := range files {
		if file.IsDir() {
			UpgradeOne(file.Name())
		}
	}
}

// UpgradeOne local repository
func UpgradeOne(repository string) error {
	termCmd := execCommand("git", "pull", "--rebase", repository)
	termCmd.Dir = config.Vipers["config"].GetString("repositories_path")

	return command.ExecuteCommand(termCmd)
}

// Remove local repository
func Remove(repository string) error {
	if repository == "/*" || repository == "/" {
		log.Errorln("Cmon, don't do that...")
	}
	termCmd := execCommand("rm", "-rf", repository)
	termCmd.Dir = config.Vipers["config"].GetString("repositories_path")

	return command.MustExecuteCommand(termCmd)
}

// Status local repository
func Status(repository string) error {
	termCmd := execCommand("git", "status")
	termCmd.Dir = config.Vipers["config"].GetString("repositories_path") + "/" + repository

	return command.ExecuteCommand(termCmd)
}

// GetGitRepositorySSHPath returns for a given repository path the full SSH path.
func GetGitRepositorySSHPath(repository string) string {
	repository = strings.TrimPrefix(strings.TrimSuffix(repository, ".git"), "https://github.com/")
	if !strings.HasPrefix(repository, "git@github.com:") {
		repository = fmt.Sprintf("git@github.com:%s.git", repository)
	}
	return repository
}
