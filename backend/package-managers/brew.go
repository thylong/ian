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

package packagemanagers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/thylong/ian/backend/command"
)

// Brew immutable instance.
var Brew = BrewPackageManager{Path: filepath.Clean("/usr/local/bin/brew"), Name: "brew"}

// BrewPackageManager is a (widely used) unofficial Mac OS package manager.
// (more: https://brew.sh/)
type BrewPackageManager struct {
	Path string
	Name string
}

// Install given Brew package.
func (brew *BrewPackageManager) Install(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(brew.Path, "install", packageName)); err != nil {
		return fmt.Errorf("Cannot %s install %s: %s", brew.Name, packageName, err)
	}
	return err
}

// Uninstall given Brew package.
func (brew *BrewPackageManager) Uninstall(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(brew.Path, "uninstall", packageName)); err != nil {
		return fmt.Errorf("Cannot %s uninstall %s: %s", brew.Name, packageName, err)
	}
	return err
}

// Cleanup all the local archives and previous versions.
func (brew *BrewPackageManager) Cleanup() (err error) {
	err = command.ExecuteCommand(execCommand(brew.Path, "cleanup"))
	command.ExecuteCommand(execCommand(brew.Path, "cask", "cleanup"))
	return err
}

// UpdateOne pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (brew *BrewPackageManager) UpdateOne(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(brew.Path, "update")); err != nil {
		return fmt.Errorf("Cannot %s update %s: %s", brew.Name, packageName, err)
	}
	return err
}

// UpgradeOne Brew packages to the last known versions.
func (brew *BrewPackageManager) UpgradeOne(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(brew.Path, "upgrade")); err != nil {
		return fmt.Errorf("Cannot %s upgrade %s: %s", brew.Name, packageName, err)
	}
	return err
}

// UpdateAll pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (brew *BrewPackageManager) UpdateAll() (err error) {
	if err = command.ExecuteCommand(execCommand(brew.Path, "update")); err != nil {
		return fmt.Errorf("Cannot %s update: %s", brew.Name, err)
	}
	return err
}

// UpgradeAll Brew packages to the last known versions.
func (brew *BrewPackageManager) UpgradeAll() (err error) {
	if err = command.ExecuteCommand(execCommand(brew.Path, "upgrade")); err != nil {
		return fmt.Errorf("Cannot %s upgrade: %s", brew.Name, err)
	}
	return err
}

// IsInstalled returns true if Brew executable is found.
func (brew *BrewPackageManager) IsInstalled() bool {
	if _, err := os.Stat(brew.Path); err != nil {
		return false
	}
	return true
}

// IsOSPackageManager returns true for Mac OS.
func (brew *BrewPackageManager) IsOSPackageManager() bool {
	return runtime.GOOS == "darwin"
}

// GetExecPath return immutable path to Brew executable.
func (brew *BrewPackageManager) GetExecPath() string {
	return brew.Path
}

// GetName return the name of the package manager.
func (brew *BrewPackageManager) GetName() string {
	return brew.Name
}

// Setup installs Brew
func (brew *BrewPackageManager) Setup() (err error) {
	resp, err := http.Get(
		"https://raw.githubusercontent.com/Homebrew/install/master/install",
	)
	if err != nil {
		return fmt.Errorf("Error: %s", err)
	}

	installScript, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	command.ExecuteInteractiveCommand(
		execCommand("/usr/bin/ruby", "-e", string(installScript)),
	)
	return nil
}
