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

	"github.com/thylong/ian/backend/command"
)

// RubyGems immutable instance.
var RubyGems = RubyGemsPackageManager{Path: filepath.Clean("/usr/local/bin/gem"), Name: "rubygems"}

// ErrRubyGemsMissingFeature is returned when triggering an unsupported feature.
var ErrRubyGemsMissingFeature = errors.New("gems is not designed to support this feature")

// RubyGemsPackageManager is a (widely used) unofficial Mac OS package manager.
// (more: https://pip.sh/)
type RubyGemsPackageManager struct {
	Path string
	Name string
}

// Install given RubyGems package.
func (pm *RubyGemsPackageManager) Install(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "install", packageName)); err != nil {
		return fmt.Errorf("Cannot %s install %s: %s", pm.Name, packageName, err)
	}
	return err
}

// Uninstall given RubyGems package.
func (pm *RubyGemsPackageManager) Uninstall(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "uninstall", packageName)); err != nil {
		return fmt.Errorf("Cannot %s uninstall %s: %s", pm.Name, packageName, err)
	}
	return err
}

// Cleanup the pip cache.
// This is done by default since pip 6.0
func (pm *RubyGemsPackageManager) Cleanup() error {
	return command.ExecuteCommand(execCommand(pm.Path, "cleanup"))
}

// UpdateOne pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (pm *RubyGemsPackageManager) UpdateOne(packageName string) (err error) {
	return ErrRubyGemsMissingFeature
}

// UpgradeOne RubyGems packages to the last known versions.
func (pm *RubyGemsPackageManager) UpgradeOne(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(pm.Path, "update", packageName)); err != nil {
		return fmt.Errorf("Cannot %s update %s: %s", pm.Name, packageName, err)
	}
	return err
}

// UpdateAll pulls last versions infos from realted repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (pm *RubyGemsPackageManager) UpdateAll() (err error) {
	return ErrRubyGemsMissingFeature
}

// UpgradeAll RubyGems packages to the last known versions.
func (pm *RubyGemsPackageManager) UpgradeAll() (err error) {
	return command.ExecuteCommand(execCommand(pm.Path, "update"))
}

// IsInstalled returns true if RubyGems executable is found.
func (pm *RubyGemsPackageManager) IsInstalled() bool {
	if _, err := os.Stat(pm.Path); err != nil {
		return false
	}
	return true
}

// IsOSPackageManager returns false because pip is never the main OS Package Manager.
func (pm *RubyGemsPackageManager) IsOSPackageManager() bool {
	return false
}

// GetExecPath return immutable path to RubyGems executable.
func (pm *RubyGemsPackageManager) GetExecPath() string {
	return pm.Path
}

// GetName return the name of the package manager.
func (pm *RubyGemsPackageManager) GetName() string {
	return pm.Name
}

// Setup installs Cask
func (pm *RubyGemsPackageManager) Setup() (err error) {
	return nil
}
