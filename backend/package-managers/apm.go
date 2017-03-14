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

// Apm immutable instance.
var Apm = ApmPackageManager{Path: "/usr/local/bin/apm", Name: "apm"}

// ApmPackageManager is the package manager for Atom text editor.
// (more: https://github.com/atom/apm)
type ApmPackageManager struct {
	Path string
	Name string
}

// Install given Apm package.
func (b ApmPackageManager) Install(packageName string) (err error) {
	err = command.ExecuteCommand(execCommand(b.Path, "install", packageName))
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	return err
}

// Uninstall given Apm package.
func (b ApmPackageManager) Uninstall(packageName string) (err error) {
	err = command.ExecuteCommand(execCommand(b.Path, "uninstall", packageName))
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	return err
}

// Cleanup all the local archives and previous versions.
func (b ApmPackageManager) Cleanup() (err error) {
	return command.ExecuteCommand(execCommand(b.Path, "clean"))
}

// UpdateOne pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (b ApmPackageManager) UpdateOne(packageName string) error {
	return command.ExecuteCommand(execCommand(b.Path, "update"))
}

// UpgradeOne Apm packages to the last known versions.
func (b ApmPackageManager) UpgradeOne(packageName string) error {
	return command.ExecuteCommand(execCommand(b.Path, "upgrade", packageName))
}

// UpdateAll pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (b ApmPackageManager) UpdateAll() error {
	return command.ExecuteCommand(execCommand(b.Path, "update"))
}

// UpgradeAll Apm packages to the last known versions.
func (b ApmPackageManager) UpgradeAll() error {
	return command.ExecuteCommand(execCommand(b.Path, "upgrade"))
}

// IsInstalled returns true if Apm executable is found.
func (b ApmPackageManager) IsInstalled() bool {
	if _, err := os.Stat(b.Path); err != nil {
		return false
	}
	return true
}

// IsOSPackageManager returns true for Mac OS.
func (b ApmPackageManager) IsOSPackageManager() bool {
	return false
}

// GetExecPath return immutable path to Apm executable.
func (b ApmPackageManager) GetExecPath() string {
	return b.Path
}

// GetName return the name of the package manager.
func (b ApmPackageManager) GetName() string {
	return b.Name
}

// Setup installs Apm
func (b ApmPackageManager) Setup() (err error) {
	return nil
}
