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
	"errors"
	"fmt"
	"os"

	"github.com/thylong/ian/backend/command"
)

// Pip immutable instance.
var Pip = PipPackageManager{Path: "/usr/local/bin/pip", Name: "pip"}

// ErrPipMissingFeature is returned when triggering an unsupported feature.
var ErrPipMissingFeature = errors.New("pip is not designed to support this feature")

// PipPackageManager is a (widely used) unofficial Mac OS package manager.
// (more: https://pip.sh/)
type PipPackageManager struct {
	Path string
	Name string
}

// Install given Pip package.
func (b PipPackageManager) Install(packageName string) (err error) {
	err = command.ExecuteCommand(execCommand(b.Path, "install", "-U", packageName))
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	return err
}

// Uninstall given Pip package.
func (b PipPackageManager) Uninstall(packageName string) (err error) {
	err = command.ExecuteCommand(execCommand(b.Path, "uninstall", "-U", packageName))
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	return err
}

// Cleanup the pip cache.
// This is done by default since pip 6.0
func (b PipPackageManager) Cleanup() (err error) {
	return nil
}

// UpdateOne pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (b PipPackageManager) UpdateOne(packageName string) (err error) {
	return nil
}

// UpgradeOne Pip packages to the last known versions.
func (b PipPackageManager) UpgradeOne(packageName string) (err error) {
	err = b.Install(packageName)
	return err
}

// UpdateAll pulls last versions infos from realted repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (b PipPackageManager) UpdateAll() (err error) {
	// TODO: Implementation
	return err
}

// UpgradeAll Pip packages to the last known versions.
func (b PipPackageManager) UpgradeAll() (err error) {
	// TODO: Implementation
	return err
}

// IsInstalled returns true if Pip executable is found.
func (b PipPackageManager) IsInstalled() bool {
	if _, err := os.Stat(b.Path); err != nil {
		return false
	}
	return true
}

// IsOSPackageManager returns false because pip is never the main OS Package Manager.
func (b PipPackageManager) IsOSPackageManager() bool {
	return false
}

// GetExecPath return immutable path to Pip executable.
func (b PipPackageManager) GetExecPath() string {
	return b.Path
}

// GetName return the name of the package manager.
func (b PipPackageManager) GetName() string {
	return b.Name
}

// Setup installs Cask
func (b PipPackageManager) Setup() (err error) {
	return nil
}
