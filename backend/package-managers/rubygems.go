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

// RubyGems immutable instance.
var RubyGems = RubyGemsPackageManager{Path: "/usr/local/bin/gem", Name: "rubygems"}

// ErrRubyGemsMissingFeature is returned when triggering an unsupported feature.
var ErrRubyGemsMissingFeature = errors.New("gems is not designed to support this feature")

// RubyGemsPackageManager is a (widely used) unofficial Mac OS package manager.
// (more: https://pip.sh/)
type RubyGemsPackageManager struct {
	Path string
	Name string
}

// Install given RubyGems package.
func (b RubyGemsPackageManager) Install(packageName string) (err error) {
	err = command.ExecuteCommand(execCommand(b.Path, "install", packageName))
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	return err
}

// Uninstall given RubyGems package.
func (b RubyGemsPackageManager) Uninstall(packageName string) (err error) {
	err = command.ExecuteCommand(execCommand(b.Path, "uninstall", packageName))
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	return err
}

// Cleanup the pip cache.
// This is done by default since pip 6.0
func (b RubyGemsPackageManager) Cleanup() error {
	return command.ExecuteCommand(execCommand(b.Path, "cleanup"))
}

// UpdateOne pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (b RubyGemsPackageManager) UpdateOne(packageName string) (err error) {
	return ErrRubyGemsMissingFeature
}

// UpgradeOne RubyGems packages to the last known versions.
func (b RubyGemsPackageManager) UpgradeOne(packageName string) error {
	return command.ExecuteCommand(execCommand(b.Path, "update", packageName))
}

// UpdateAll pulls last versions infos from realted repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (b RubyGemsPackageManager) UpdateAll() error {
	return ErrRubyGemsMissingFeature
}

// UpgradeAll RubyGems packages to the last known versions.
func (b RubyGemsPackageManager) UpgradeAll() error {
	return command.ExecuteCommand(execCommand(b.Path, "update"))
}

// IsInstalled returns true if RubyGems executable is found.
func (b RubyGemsPackageManager) IsInstalled() bool {
	if _, err := os.Stat(b.Path); err != nil {
		return false
	}
	return true
}

// IsOSPackageManager returns false because pip is never the main OS Package Manager.
func (b RubyGemsPackageManager) IsOSPackageManager() bool {
	return false
}

// GetExecPath return immutable path to RubyGems executable.
func (b RubyGemsPackageManager) GetExecPath() string {
	return b.Path
}

// GetName return the name of the package manager.
func (b RubyGemsPackageManager) GetName() string {
	return b.Name
}

// Setup installs Cask
func (b RubyGemsPackageManager) Setup() (err error) {
	return nil
}
