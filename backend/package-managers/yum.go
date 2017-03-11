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

// Yum immutable instance.
var Yum = YumPackageManager{Path: "/usr/bin/yum", Name: "yum"}

// YumPackageManager is the official Debian (and associated distributions) package manager.
// (more: https://wiki.debian.org/Yum)
type YumPackageManager struct {
	Path string
	Name string
}

// Install given Yum package.
func (b YumPackageManager) Install(packageName string) (err error) {
	err = command.ExecuteCommand(execCommand(b.Path, "install", packageName))
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	return err
}

// Uninstall given Yum package.
func (b YumPackageManager) Uninstall(packageName string) (err error) {
	err = command.ExecuteCommand(execCommand(b.Path, "erase", packageName))
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	return err
}

// Cleanup all the local archives and previous versions.
func (b YumPackageManager) Cleanup() error {
	return command.ExecuteCommand(execCommand(b.Path, "autoremove"))
}

// UpdateOne pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (b YumPackageManager) UpdateOne(packageName string) error {
	return command.ExecuteCommand(execCommand(b.Path, "update", packageName))
}

// UpgradeOne Yum packages to the last known versions.
func (b YumPackageManager) UpgradeOne(packageName string) error {
	return command.ExecuteCommand(execCommand(b.Path, "upgrade", packageName))
}

// UpdateAll pulls last versions infos from realted repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (b YumPackageManager) UpdateAll() error {
	return command.ExecuteCommand(execCommand(b.Path, "update"))
}

// UpgradeAll Yum packages to the last known versions.
func (b YumPackageManager) UpgradeAll() error {
	return command.ExecuteCommand(execCommand(b.Path, "upgrade"))
}

// IsInstalled returns true if Yum executable is found.
func (b YumPackageManager) IsInstalled() bool {
	if fileInfo, err := os.Stat(b.Path); err != nil || fileInfo.Mode() == os.ModeSymlink {
		return false
	}
	return true
}

// IsOSPackageManager returns true for Mac OS.
func (b YumPackageManager) IsOSPackageManager() bool {
	return b.IsInstalled()
}

// GetExecPath return immutable path to Yum executable.
func (b YumPackageManager) GetExecPath() string {
	return b.Path
}

// GetName return the name of the package manager.
func (b YumPackageManager) GetName() string {
	return b.Name
}

// Setup does nothing (yum comes by default in Linux distributions)
func (b YumPackageManager) Setup() (err error) {
	return nil
}
