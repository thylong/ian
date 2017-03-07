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
	"os/exec"

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
func (b CaskPackageManager) Install(packageName string) (err error) {
	err = command.ExecuteCommand(exec.Command(b.Path, "cask", "install", packageName))
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	return err
}

// Uninstall given Cask package.
func (b CaskPackageManager) Uninstall(packageName string) (err error) {
	err = command.ExecuteCommand(exec.Command(b.Path, "cask", "uninstall", packageName))
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	return err
}

// Cleanup all the local archives and previous versions.
func (b CaskPackageManager) Cleanup() (err error) {
	// Cleanup brew
	err = command.ExecuteCommand(exec.Command(b.Path, "cask", "cleanup"))
	// Cleanup cask
	command.ExecuteCommand(exec.Command("/usr/local/bin/brew", "cask", "cleanup"))
	return err
}

// UpdateOne pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (b CaskPackageManager) UpdateOne(packageName string) (err error) {
	err = command.ExecuteCommand(exec.Command(b.Path, "cask", "update"))
	return err
}

// UpgradeOne Cask packages to the last known versions.
func (b CaskPackageManager) UpgradeOne(packageName string) (err error) {
	err = command.ExecuteCommand(exec.Command(b.Path, "cask", "upgrade", packageName))
	return err
}

// UpdateAll pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (b CaskPackageManager) UpdateAll() (err error) {
	err = command.ExecuteCommand(exec.Command(b.Path, "cask", "update"))
	return err
}

// UpgradeAll Cask packages to the last known versions.
func (b CaskPackageManager) UpgradeAll() (err error) {
	err = command.ExecuteCommand(exec.Command(b.Path, "cask", "upgrade"))
	return err
}

// IsInstalled returns true if Cask executable is found.
func (b CaskPackageManager) IsInstalled() bool {
	if _, err := os.Stat("/usr/local/bin/cask"); err != nil {
		return false
	}
	return true
}

// IsOSPackageManager returns true for Mac OS.
func (b CaskPackageManager) IsOSPackageManager() bool {
	return false
}

// GetExecPath return immutable path to Cask executable.
func (b CaskPackageManager) GetExecPath() string {
	return b.Path
}

// GetName return the name of the package manager.
func (b CaskPackageManager) GetName() string {
	return b.Name
}

// Setup installs Cask
func (b CaskPackageManager) Setup() (err error) {
	fmt.Print("Installing cask...")
	if _, err := os.Stat("/usr/local/bin/cask"); err != nil {
		err = command.ExecuteCommand(exec.Command("brew", "tap", "caskroom/cask"))
		return err
	}
	fmt.Print("cask already installed, skipping...")
	return nil
}
