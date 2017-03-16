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
func (pm ApmPackageManager) Install(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "install", packageName)); err != nil {
		return fmt.Errorf("Cannot %s install %s: %s", pm.Name, packageName, err)
	}
	return err
}

// Uninstall given Apm package.
func (pm ApmPackageManager) Uninstall(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "uninstall", packageName)); err != nil {
		return fmt.Errorf("Cannot %s uninstall %s: %s", pm.Name, packageName, err)
	}
	return err
}

// Cleanup all the local archives and previous versions.
func (pm ApmPackageManager) Cleanup() (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "clean")); err != nil {
		return fmt.Errorf("Cannot %s clean: %s", pm.Name, err)
	}
	return err
}

// UpdateOne pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (pm ApmPackageManager) UpdateOne(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "update")); err != nil {
		return fmt.Errorf("Cannot %s update %s: %s", pm.Name, packageName, err)
	}
	return err
}

// UpgradeOne Apm packages to the last known versions.
func (pm ApmPackageManager) UpgradeOne(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "upgrade", packageName)); err != nil {
		return fmt.Errorf("Cannot %s upgrade %s: %s", pm.Name, packageName, err)
	}
	return err
}

// UpdateAll pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (pm ApmPackageManager) UpdateAll() (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "update")); err != nil {
		return fmt.Errorf("Cannot %s update: %s", pm.Name, err)
	}
	return err
}

// UpgradeAll Apm packages to the last known versions.
func (pm ApmPackageManager) UpgradeAll() (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "upgrade")); err != nil {
		return fmt.Errorf("Cannot %s upgrade: %s", pm.Name, err)
	}
	return err
}

// IsInstalled returns true if Apm executable is found.
func (pm ApmPackageManager) IsInstalled() bool {
	if _, err := os.Stat(pm.Path); err != nil {
		return false
	}
	return true
}

// IsOSPackageManager returns true for Mac OS.
func (pm ApmPackageManager) IsOSPackageManager() bool {
	return false
}

// GetExecPath return immutable path to Apm executable.
func (pm ApmPackageManager) GetExecPath() string {
	return pm.Path
}

// GetName return the name of the package manager.
func (pm ApmPackageManager) GetName() string {
	return pm.Name
}

// Setup installs Apm
func (pm ApmPackageManager) Setup() (err error) {
	return nil
}
