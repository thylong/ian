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

// Npm immutable instance.
var Npm = NpmPackageManager{Path: "/usr/local/bin/npm", Name: "npm"}

// ErrNPMMissingFeature is returned when triggering an unsupported feature.
var ErrNPMMissingFeature = errors.New("npm is not designed to support this feature")

// NpmPackageManager is a (widely used) unofficial Mac OS package manager.
// (more: https://npm.sh/)
type NpmPackageManager struct {
	Path string
	Name string
}

// Install given Npm package.
func (pm *NpmPackageManager) Install(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "install", "-g", packageName)); err != nil {
		return fmt.Errorf("Cannot %s install %s: %s", pm.Name, packageName, err)
	}
	return err
}

// Uninstall given Npm package.
func (pm *NpmPackageManager) Uninstall(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "uninstall", "-g", packageName)); err != nil {
		return fmt.Errorf("Cannot %s install %s: %s", pm.Name, packageName, err)
	}
	return err
}

// Cleanup the npm cache.
func (pm *NpmPackageManager) Cleanup() error {
	return command.ExecuteCommand(execCommand(pm.Path, "cache", "clean"))
}

// UpdateOne pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (pm *NpmPackageManager) UpdateOne(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "update", packageName)); err != nil {
		return fmt.Errorf("Cannot %s update %s: %s", pm.Name, packageName, err)
	}
	return err
}

// UpgradeOne Npm packages to the last known versions.
func (pm *NpmPackageManager) UpgradeOne(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "upgrade", packageName)); err != nil {
		return fmt.Errorf("Cannot %s upgrade %s: %s", pm.Name, packageName, err)
	}
	return err
}

// UpdateAll does nothing (out of making NPM satisfying PackageManager interface).
func (pm *NpmPackageManager) UpdateAll() error {
	return ErrNPMMissingFeature
}

// UpgradeAll Npm packages to the last known versions.
func (pm *NpmPackageManager) UpgradeAll() (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "update", "-g")); err != nil {
		return fmt.Errorf("Cannot %s update: %s", pm.Name, err)
	}
	return err
}

// IsInstalled returns true if Npm executable is found.
func (pm *NpmPackageManager) IsInstalled() bool {
	if _, err := os.Stat(pm.Path); err != nil {
		return false
	}
	return true
}

// IsOSPackageManager returns false because npm is never the main OS Package Manager.
func (pm *NpmPackageManager) IsOSPackageManager() bool {
	return false
}

// GetExecPath return immutable path to Npm executable.
func (pm *NpmPackageManager) GetExecPath() string {
	return pm.Path
}

// GetName return the name of the package manager.
func (pm *NpmPackageManager) GetName() string {
	return pm.Name
}

// Setup installs Cask
func (pm *NpmPackageManager) Setup() (err error) {
	return nil
}
