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
	"runtime"

	"github.com/thylong/ian/backend/command"
)

// Apt immutable instance.
var Apt = AptPackageManager{Path: "/usr/bin/apt-get", Name: "apt"}

// ErrAptMissingFeature is returned when triggering an unsupported feature.
var ErrAptMissingFeature = errors.New("apt is not designed to support this feature")

// AptPackageManager is the official Debian (and associated distributions) package manager.
// (more: https://wiki.debian.org/Apt)
type AptPackageManager struct {
	Path string
	Name string
}

// Install given Apt package.
func (pm AptPackageManager) Install(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "install", packageName)); err != nil {
		return fmt.Errorf("Cannot %s install: %s", pm.Name, err)
	}
	return err
}

// Uninstall given Apt package.
func (pm AptPackageManager) Uninstall(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "remove", packageName)); err != nil {
		return fmt.Errorf("Cannot %s remove: %s", pm.Name, err)
	}
	return err
}

// Cleanup all the local archives and previous versions.
func (pm AptPackageManager) Cleanup() (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "autoremove")); err != nil {
		return fmt.Errorf("Cannot %s autoremove: %s", pm.Name, err)
	}
	return err
}

// UpdateOne pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (pm AptPackageManager) UpdateOne(packageName string) (err error) {
	return ErrAptMissingFeature
}

// UpgradeOne Npm packages to the last known versions.
func (pm AptPackageManager) UpgradeOne(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "upgrade", packageName)); err != nil {
		return fmt.Errorf("Cannot %s upgrade: %s", pm.Name, err)
	}
	return err
}

// UpdateAll pulls last versions infos from realted repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (pm AptPackageManager) UpdateAll() (err error) {
	return command.ExecuteCommand(execCommand(pm.Path, "update"))
}

// UpgradeAll Apt packages to the last known versions.
func (pm AptPackageManager) UpgradeAll() (err error) {
	return command.ExecuteCommand(execCommand(pm.Path, "full-upgrade"))
}

// IsInstalled returns true if Apt executable is found.
func (pm AptPackageManager) IsInstalled() bool {
	if _, err := os.Stat(pm.Path); err != nil {
		return false
	}
	return true
}

// IsOSPackageManager returns true for Mac OS.
func (pm AptPackageManager) IsOSPackageManager() bool {
	return pm.IsInstalled() && runtime.GOOS == "linux"
}

// GetExecPath return immutable path to Apt executable.
func (pm AptPackageManager) GetExecPath() string {
	return pm.Path
}

// GetName return the name of the package manager.
func (pm AptPackageManager) GetName() string {
	return pm.Name
}

// Setup does nothing (apt comes by default in Linux distributions)
func (pm AptPackageManager) Setup() (err error) {
	return nil
}
