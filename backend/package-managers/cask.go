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
	"os"
	"path/filepath"

	"github.com/thylong/ian/backend/command"
)

// Cask immutable instance.
var Cask = CaskPackageManager{Path: filepath.Clean("/usr/local/bin/brew"), Name: "cask"}

var caskPath = filepath.Clean("/usr/local/bin/cask")

// CaskPackageManager is an extension of Brew Mac OS package manager.
// (more: https://caskroom.github.io/)
type CaskPackageManager struct {
	Path string
	Name string
}

// Install given Cask package.
func (cask *CaskPackageManager) Install(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(cask.Path, "cask", "install", packageName)); err != nil {
		return fmt.Errorf("Cannot %s install %s: %s", cask.Name, packageName, err)
	}
	return err
}

// Uninstall given Cask package.
func (cask *CaskPackageManager) Uninstall(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(cask.Path, "cask", "uninstall", packageName)); err != nil {
		return fmt.Errorf("Cannot %s uninstall %s: %s", cask.Name, packageName, err)
	}
	return err
}

// Cleanup all the local archives and previous versions.
func (cask *CaskPackageManager) Cleanup() error {
	return command.ExecuteCommand(execCommand(cask.Path, "cask", "cleanup"))
}

// UpdateOne pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (cask *CaskPackageManager) UpdateOne(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(cask.Path, "cask", "update")); err != nil {
		return fmt.Errorf("Cannot %s update %s: %s", cask.Name, packageName, err)
	}
	return err
}

// UpgradeOne Cask packages to the last known versions.
func (cask *CaskPackageManager) UpgradeOne(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(cask.Path, "cask", "upgrade")); err != nil {
		return fmt.Errorf("Cannot %s upgrade %s: %s", cask.Name, packageName, err)
	}
	return err
}

// UpdateAll pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (cask *CaskPackageManager) UpdateAll() (err error) {
	if err = command.ExecuteCommand(execCommand(cask.Path, "cask", "update")); err != nil {
		return fmt.Errorf("Cannot %s update: %s", cask.Name, err)
	}
	return err
}

// UpgradeAll Cask packages to the last known versions.
func (cask *CaskPackageManager) UpgradeAll() (err error) {
	if err = command.ExecuteCommand(execCommand(cask.Path, "cask", "upgrade")); err != nil {
		return fmt.Errorf("Cannot %s upgrade: %s", cask.Name, err)
	}
	return err
}

// IsInstalled returns true if Cask executable is found.
func (cask *CaskPackageManager) IsInstalled() bool {
	if _, err := os.Stat(caskPath); err != nil {
		return false
	}
	return true
}

// IsOSPackageManager returns true for Mac OS.
func (cask *CaskPackageManager) IsOSPackageManager() bool {
	return false
}

// GetExecPath return immutable path to Cask executable.
func (cask *CaskPackageManager) GetExecPath() string {
	return cask.Path
}

// GetName return the name of the package manager.
func (cask *CaskPackageManager) GetName() string {
	return cask.Name
}

// Setup installs Cask
func (cask *CaskPackageManager) Setup() (err error) {
	fmt.Print("Installing cask...")
	if _, err := os.Stat(caskPath); err != nil {
		return command.ExecuteCommand(execCommand("brew", "tap", "caskroom/cask"))
	}
	fmt.Print("cask already installed, skipping...")
	return nil
}
