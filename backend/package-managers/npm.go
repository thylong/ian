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
func (b NpmPackageManager) Install(packageName string) (err error) {
	err = command.ExecuteCommand(execCommand(b.Path, "install", "-g", packageName))
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	return err
}

// Uninstall given Npm package.
func (b NpmPackageManager) Uninstall(packageName string) (err error) {
	err = command.ExecuteCommand(execCommand(b.Path, "uninstall", "-g", packageName))
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	return err
}

// Cleanup the npm cache.
func (b NpmPackageManager) Cleanup() error {
	return command.ExecuteCommand(execCommand(b.Path, "cache", "clean"))
}

// UpdateOne pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (b NpmPackageManager) UpdateOne(packageName string) error {
	return command.ExecuteCommand(execCommand(b.Path, "update", packageName))
}

// UpgradeOne Npm packages to the last known versions.
func (b NpmPackageManager) UpgradeOne(packageName string) error {
	return command.ExecuteCommand(execCommand(b.Path, "upgrade", packageName))
}

// UpdateAll does nothing (out of making NPM satisfying PackageManager interface).
func (b NpmPackageManager) UpdateAll() error {
	return ErrNPMMissingFeature
}

// UpgradeAll Npm packages to the last known versions.
func (b NpmPackageManager) UpgradeAll() error {
	return command.ExecuteCommand(execCommand(b.Path, "update", "-g"))
}

// IsInstalled returns true if Npm executable is found.
func (b NpmPackageManager) IsInstalled() bool {
	if _, err := os.Stat(b.Path); err != nil {
		return false
	}
	return true
}

// IsOSPackageManager returns false because npm is never the main OS Package Manager.
func (b NpmPackageManager) IsOSPackageManager() bool {
	return false
}

// GetExecPath return immutable path to Npm executable.
func (b NpmPackageManager) GetExecPath() string {
	return b.Path
}

// GetName return the name of the package manager.
func (b NpmPackageManager) GetName() string {
	return b.Name
}

// Setup installs Cask
func (b NpmPackageManager) Setup() (err error) {
	return nil
}
