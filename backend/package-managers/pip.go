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
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/thylong/ian/backend/command"
)

// Pip immutable instance.
var Pip = PipPackageManager{Path: GetDefaultPath(), Name: "pip"}

// ErrPipMissingFeature is returned when triggering an unsupported feature.
var ErrPipMissingFeature = errors.New("pip is not designed to support this feature")

// PipPackageManager is a (widely used) unofficial Mac OS package manager.
// (more: https://pip.sh/)
type PipPackageManager struct {
	Path string
	Name string
}

// Install given Pip package.
func (pip *PipPackageManager) Install(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pip.Path, "install", "-U", packageName)); err != nil {
		return fmt.Errorf("Cannot %s install %s: %s", pip.Name, packageName, err)
	}
	return err
}

// Uninstall given Pip package.
func (pip *PipPackageManager) Uninstall(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pip.Path, "uninstall", "-U", packageName)); err != nil {
		return fmt.Errorf("Cannot %s uninstall %s: %s", pip.Name, packageName, err)
	}
	return err
}

// Cleanup the pip cache.
// This is done by default since pip 6.0
func (pip *PipPackageManager) Cleanup() error {
	return ErrPipMissingFeature
}

// UpdateOne pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (pip *PipPackageManager) UpdateOne(packageName string) (err error) {
	return ErrPipMissingFeature
}

// UpgradeOne Pip packages to the last known versions.
func (pip *PipPackageManager) UpgradeOne(packageName string) (err error) {
	return pip.Install(packageName)
}

// UpdateAll pulls last versions infos from realted repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (pip *PipPackageManager) UpdateAll() (err error) {
	// TODO: Implementation
	return ErrPipMissingFeature
}

// UpgradeAll Pip packages to the last known versions.
func (pip *PipPackageManager) UpgradeAll() (err error) {
	// TODO: Implementation
	return ErrPipMissingFeature
}

// IsInstalled returns true if Pip executable is found.
func (pip *PipPackageManager) IsInstalled() bool {
	if _, err := os.Stat(pip.Path); err != nil {
		return false
	}
	return true
}

// IsOSPackageManager returns false because pip is never the main OS Package Manager.
func (pip *PipPackageManager) IsOSPackageManager() bool {
	return false
}

// GetExecPath return immutable path to Pip executable.
func (pip *PipPackageManager) GetExecPath() string {
	return pip.Path
}

// GetDefaultPath return default path to Pip executable.
func GetDefaultPath() string {
	defaultPath := filepath.Clean("/usr/local/bin/pip")
	if _, err := os.Stat(defaultPath); err != nil {
		if _, err := os.Stat(fmt.Sprintf("%s%s", defaultPath, "2")); err != nil {
			return defaultPath
		}
		return fmt.Sprintf("%s%s", defaultPath, "2")
	}
	return defaultPath
}

// GetName return the name of the package manager.
func (pip *PipPackageManager) GetName() string {
	return pip.Name
}

// Setup installs Cask
func (pip *PipPackageManager) Setup() (err error) {
	return nil
}
