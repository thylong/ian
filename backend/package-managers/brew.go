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

package packagemanagers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"

	"github.com/thylong/ian/backend/command"
)

// Brew immutable instance.
var Brew = BrewPackageManager{Path: "/usr/local/bin/brew", Name: "brew"}

// BrewPackageManager is a (widely used) unofficial Mac OS package manager.
// (more: https://brew.sh/)
type BrewPackageManager struct {
	Path string
	Name string
}

// Install given Brew package.
func (b BrewPackageManager) Install(packageName string) (err error) {
	err = command.ExecuteCommand(execCommand(b.Path, "install", packageName))
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	return err
}

// Uninstall given Brew package.
func (b BrewPackageManager) Uninstall(packageName string) (err error) {
	err = command.ExecuteCommand(execCommand(b.Path, "uninstall", packageName))
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	return err
}

// Cleanup all the local archives and previous versions.
func (b BrewPackageManager) Cleanup() (err error) {
	err = command.ExecuteCommand(execCommand(b.Path, "cleanup"))
	command.ExecuteCommand(execCommand(b.Path, "cask", "cleanup"))
	return err
}

// UpdateOne pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (b BrewPackageManager) UpdateOne(packageName string) error {
	return command.ExecuteCommand(execCommand(b.Path, "update"))
}

// UpgradeOne Brew packages to the last known versions.
func (b BrewPackageManager) UpgradeOne(packageName string) error {
	return command.ExecuteCommand(execCommand(b.Path, "upgrade", packageName))
}

// UpdateAll pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (b BrewPackageManager) UpdateAll() error {
	return command.ExecuteCommand(execCommand(b.Path, "update"))
}

// UpgradeAll Brew packages to the last known versions.
func (b BrewPackageManager) UpgradeAll() error {
	return command.ExecuteCommand(execCommand(b.Path, "upgrade"))
}

// IsInstalled returns true if Brew executable is found.
func (b BrewPackageManager) IsInstalled() bool {
	if _, err := os.Stat(b.Path); err != nil {
		return false
	}
	return true
}

// IsOSPackageManager returns true for Mac OS.
func (b BrewPackageManager) IsOSPackageManager() bool {
	return b.IsInstalled() && runtime.GOOS == "darwin"
}

// GetExecPath return immutable path to Brew executable.
func (b BrewPackageManager) GetExecPath() string {
	return b.Path
}

// GetName return the name of the package manager.
func (b BrewPackageManager) GetName() string {
	return b.Name
}

// Setup installs Brew
func (b BrewPackageManager) Setup() (err error) {
	resp, err := http.Get(
		"https://raw.githubusercontent.com/Homebrew/install/master/install",
	)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	installScript, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	installCmd := execCommand("/usr/bin/ruby", "-e", string(installScript))
	installCmd.Stdout = os.Stdout
	installCmd.Stdin = os.Stdin
	installCmd.Stderr = os.Stderr
	return installCmd.Run()
}
