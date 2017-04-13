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
	"runtime"

	"github.com/thylong/ian/backend/command"
)

// Yum immutable instance.
var Yum = YumPackageManager{Path: "/usr/bin/yum", Name: "yum"}

// YumPackageManager is the official Debian (and associated distributions) package manager.
// (more: https://wiki.debian.org/Yum)
type YumPackageManager struct {
	Path string
	Name string
}

// Install given Yum package.
func (pm *YumPackageManager) Install(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "install", packageName)); err != nil {
		return fmt.Errorf("Cannot %s install %s: %s", pm.Name, packageName, err)
	}
	return err
}

// Uninstall given Yum package.
func (pm *YumPackageManager) Uninstall(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "erase", packageName)); err != nil {
		return fmt.Errorf("Cannot %s erase %s: %s", pm.Name, packageName, err)
	}
	return err
}

// Cleanup all the local archives and previous versions.
func (pm *YumPackageManager) Cleanup() error {
	return command.ExecuteCommand(execCommand(pm.Path, "autoremove"))
}

// UpdateOne pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (pm *YumPackageManager) UpdateOne(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "update", packageName)); err != nil {
		return fmt.Errorf("Cannot %s update %s: %s", pm.Name, packageName, err)
	}
	return err
}

// UpgradeOne Yum packages to the last known versions.
func (pm *YumPackageManager) UpgradeOne(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "upgrade", packageName)); err != nil {
		return fmt.Errorf("Cannot %s upgrade %s: %s", pm.Name, packageName, err)
	}
	return err
}

// UpdateAll pulls last versions infos from realted repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (pm *YumPackageManager) UpdateAll() (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "update")); err != nil {
		return fmt.Errorf("Cannot %s update: %s", pm.Name, err)
	}
	return err
}

// UpgradeAll Yum packages to the last known versions.
func (pm *YumPackageManager) UpgradeAll() (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "upgrade")); err != nil {
		return fmt.Errorf("Cannot %s upgrade: %s", pm.Name, err)
	}
	return err
}

// IsInstalled returns true if Yum executable is found.
func (pm *YumPackageManager) IsInstalled() bool {
	if fileInfo, err := os.Stat(pm.Path); err != nil || fileInfo.Mode() == os.ModeSymlink {
		return false
	}
	return true
}

// IsOSPackageManager returns true for Mac OS.
func (pm *YumPackageManager) IsOSPackageManager() bool {
	return pm.IsInstalled() && runtime.GOOS == "linux"
}

// GetExecPath return immutable path to Yum executable.
func (pm *YumPackageManager) GetExecPath() string {
	return pm.Path
}

// GetName return the name of the package manager.
func (pm *YumPackageManager) GetName() string {
	return pm.Name
}

// Setup does nothing (yum comes by default in Linux distributions)
func (pm *YumPackageManager) Setup() (err error) {
	return nil
}
