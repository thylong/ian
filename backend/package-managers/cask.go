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
	"os"

	"github.com/thylong/ian/backend/command"
)

// Cask immutable instance.
var Cask = CaskPackageManager{Path: "/usr/local/bin/brew", Name: "cask"}

// CaskPackageManager is an extension of Brew Mac OS package manager.
// (more: https://caskroom.github.io/)
type CaskPackageManager struct {
	Path string
	Name string
}

// Install given Cask package.
func (pm *CaskPackageManager) Install(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "cask", "install", packageName)); err != nil {
		return fmt.Errorf("Cannot %s install %s: %s", pm.Name, packageName, err)
	}
	return err
}

// Uninstall given Cask package.
func (pm *CaskPackageManager) Uninstall(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "cask", "uninstall", packageName)); err != nil {
		return fmt.Errorf("Cannot %s uninstall %s: %s", pm.Name, packageName, err)
	}
	return err
}

// Cleanup all the local archives and previous versions.
func (pm *CaskPackageManager) Cleanup() error {
	return command.ExecuteCommand(execCommand(pm.Path, "cask", "cleanup"))
}

// UpdateOne pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (pm *CaskPackageManager) UpdateOne(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "cask", "update")); err != nil {
		return fmt.Errorf("Cannot %s update %s: %s", pm.Name, packageName, err)
	}
	return err
}

// UpgradeOne Cask packages to the last known versions.
func (pm *CaskPackageManager) UpgradeOne(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "cask", "upgrade")); err != nil {
		return fmt.Errorf("Cannot %s upgrade %s: %s", pm.Name, packageName, err)
	}
	return err
}

// UpdateAll pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (pm *CaskPackageManager) UpdateAll() (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "cask", "update")); err != nil {
		return fmt.Errorf("Cannot %s update: %s", pm.Name, err)
	}
	return err
}

// UpgradeAll Cask packages to the last known versions.
func (pm *CaskPackageManager) UpgradeAll() (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "cask", "upgrade")); err != nil {
		return fmt.Errorf("Cannot %s upgrade: %s", pm.Name, err)
	}
	return err
}

// IsInstalled returns true if Cask executable is found.
func (pm *CaskPackageManager) IsInstalled() bool {
	if _, err := os.Stat("/usr/local/bin/cask"); err != nil {
		return false
	}
	return true
}

// IsOSPackageManager returns true for Mac OS.
func (pm *CaskPackageManager) IsOSPackageManager() bool {
	return false
}

// GetExecPath return immutable path to Cask executable.
func (pm *CaskPackageManager) GetExecPath() string {
	return pm.Path
}

// GetName return the name of the package manager.
func (pm *CaskPackageManager) GetName() string {
	return pm.Name
}

// Setup installs Cask
func (pm *CaskPackageManager) Setup() (err error) {
	fmt.Print("Installing cask...")
	if _, err := os.Stat("/usr/local/bin/cask"); err != nil {
		return command.ExecuteCommand(execCommand("brew", "tap", "caskroom/cask"))
	}
	fmt.Print("cask already installed, skipping...")
	return nil
}
