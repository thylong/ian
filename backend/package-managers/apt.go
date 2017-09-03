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
	"path/filepath"
	"runtime"

	"github.com/thylong/ian/backend/command"
)

// Apt immutable instance.
var Apt = AptPackageManager{Path: filepath.Clean("/usr/bin/apt-get"), Name: "apt"}

// ErrAptMissingFeature is returned when triggering an unsupported feature.
var ErrAptMissingFeature = errors.New("apt is not designed to support this feature")

// AptPackageManager is the official Debian (and associated distributions) package manager.
// (more: https://wiki.debian.org/Apt)
type AptPackageManager struct {
	Path string
	Name string
}

// Install given Apt package.
func (apt *AptPackageManager) Install(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(apt.Path, "install", packageName)); err != nil {
		return fmt.Errorf("Cannot %s install: %s", apt.Name, err)
	}
	return err
}

// Uninstall given Apt package.
func (apt *AptPackageManager) Uninstall(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(apt.Path, "remove", packageName)); err != nil {
		return fmt.Errorf("Cannot %s remove: %s", apt.Name, err)
	}
	return err
}

// Cleanup all the local archives and previous versions.
func (apt *AptPackageManager) Cleanup() (err error) {
	if err = command.ExecuteCommand(execCommand(apt.Path, "autoremove")); err != nil {
		return fmt.Errorf("Cannot %s autoremove: %s", apt.Name, err)
	}
	return err
}

// UpdateOne pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (apt *AptPackageManager) UpdateOne(packageName string) (err error) {
	return ErrAptMissingFeature
}

// UpgradeOne Npm packages to the last known versions.
func (apt *AptPackageManager) UpgradeOne(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(apt.Path, "upgrade", packageName)); err != nil {
		return fmt.Errorf("Cannot %s upgrade: %s", apt.Name, err)
	}
	return err
}

// UpdateAll pulls last versions infos from realted repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (apt *AptPackageManager) UpdateAll() (err error) {
	return command.ExecuteCommand(execCommand(apt.Path, "update"))
}

// UpgradeAll Apt packages to the last known versions.
func (apt *AptPackageManager) UpgradeAll() (err error) {
	return command.ExecuteCommand(execCommand(apt.Path, "full-upgrade"))
}

// IsInstalled returns true if Apt executable is found.
func (apt *AptPackageManager) IsInstalled() bool {
	if _, err := os.Stat(apt.Path); err != nil {
		return false
	}
	return true
}

// IsOSPackageManager returns true for Mac OS.
func (apt *AptPackageManager) IsOSPackageManager() bool {
	return apt.IsInstalled() && runtime.GOOS == "linux"
}

// GetExecPath return immutable path to Apt executable.
func (apt *AptPackageManager) GetExecPath() string {
	return apt.Path
}

// GetName return the name of the package manager.
func (apt *AptPackageManager) GetName() string {
	return apt.Name
}

// Setup does nothing (apt comes by default in Linux distributions)
func (apt *AptPackageManager) Setup() (err error) {
	return nil
}
