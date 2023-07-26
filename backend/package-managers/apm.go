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

// Apm immutable instance.
var Apm = ApmPackageManager{Path: filepath.Clean("/usr/local/bin/apm"), Name: "apm"}

// ApmPackageManager is the package manager for Atom text editor.
// (more: https://github.com/atom/apm)
type ApmPackageManager struct {
	Path string
	Name string
}

// Install given Apm package.
func (apm *ApmPackageManager) Install(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(apm.Path, "install", packageName)); err != nil {
		return fmt.Errorf("Cannot %s install %s: %s", apm.Name, packageName, err)
	}
	return err
}

// Uninstall given Apm package.
func (apm *ApmPackageManager) Uninstall(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(apm.Path, "uninstall", packageName)); err != nil {
		return fmt.Errorf("Cannot %s uninstall %s: %s", apm.Name, packageName, err)
	}
	return err
}

// Cleanup all the local archives and previous versions.
func (apm *ApmPackageManager) Cleanup() (err error) {
	if err = command.ExecuteCommand(execCommand(apm.Path, "clean")); err != nil {
		return fmt.Errorf("Cannot %s clean: %s", apm.Name, err)
	}
	return err
}

// UpdateOne pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (apm *ApmPackageManager) UpdateOne(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(apm.Path, "update", "--confirm=false")); err != nil {
		return fmt.Errorf("Cannot %s update %s: %s", apm.Name, packageName, err)
	}
	return err
}

// UpgradeOne Apm packages to the last known versions.
func (apm *ApmPackageManager) UpgradeOne(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(apm.Path, "upgrade", "--confirm=false", packageName)); err != nil {
		return fmt.Errorf("Cannot %s upgrade %s: %s", apm.Name, packageName, err)
	}
	return err
}

// UpdateAll pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (apm *ApmPackageManager) UpdateAll() (err error) {
	if err = command.ExecuteCommand(execCommand(apm.Path, "update", "--confirm=false")); err != nil {
		return fmt.Errorf("Cannot %s update: %s", apm.Name, err)
	}
	return err
}

// UpgradeAll Apm packages to the last known versions.
func (apm *ApmPackageManager) UpgradeAll() (err error) {
	if err = command.ExecuteCommand(execCommand(apm.Path, "upgrade", "--confirm=false")); err != nil {
		return fmt.Errorf("Cannot %s upgrade: %s", apm.Name, err)
	}
	return err
}

// IsInstalled returns true if Apm executable is found.
func (apm *ApmPackageManager) IsInstalled() bool {
	if _, err := os.Stat(apm.Path); err != nil {
		return false
	}
	return true
}

// IsOSPackageManager returns true for Mac OS.
func (apm *ApmPackageManager) IsOSPackageManager() bool {
	return false
}

// GetExecPath return immutable path to Apm executable.
func (apm *ApmPackageManager) GetExecPath() string {
	return apm.Path
}

// GetName return the name of the package manager.
func (apm *ApmPackageManager) GetName() string {
	return apm.Name
}

// Setup installs Apm
func (apm *ApmPackageManager) Setup() (err error) {
	return nil
}
