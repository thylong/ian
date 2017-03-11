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
	"os/exec"

	"github.com/thylong/ian/backend/command"
	"github.com/thylong/ian/backend/config"
)

// PackageManager handles standard interactions with all Package Managers.
type PackageManager interface {
	Install(packageName string) error
	Uninstall(packageName string) error
	Cleanup() error
	UpdateOne(string) error
	UpgradeOne(string) error
	UpdateAll() error
	UpgradeAll() error
	IsInstalled() bool
	IsOSPackageManager() bool
	GetExecPath() string
	GetName() string
	Setup() error
}

// SupportedPackageManagers contains all the currently supported package managers.
var SupportedPackageManagers = make(map[string]PackageManager)

var execCommand = exec.Command

func init() {
	SupportedPackageManagers["brew"] = Brew
	SupportedPackageManagers["cask"] = Cask
	SupportedPackageManagers["pip"] = Pip
	SupportedPackageManagers["npm"] = Npm
	SupportedPackageManagers["apt"] = Apt
	SupportedPackageManagers["yum"] = Yum
	SupportedPackageManagers["rubygems"] = RubyGems
}

// GetOSPackageManager returns the main Package Manager of the current OS.
// As only MacOS is supported for now, it returns a Brew instance.
func GetOSPackageManager() PackageManager {
	for name, packageManager := range SupportedPackageManagers {
		if name != "cask" && packageManager.IsOSPackageManager() {
			return packageManager
		}
	}
	fmt.Println("OS not supported yet.")
	return Brew
}

// GetPackageManager returns the corresponding PackageManager otherwise default
// to OS package manager.
func GetPackageManager(PackageManagerFlag string) PackageManager {
	packageManager, ok := SupportedPackageManagers[PackageManagerFlag]

	if ok {
		return packageManager
	}
	return GetOSPackageManager()
}

// SearchOnPackageManagers returns infos on packages found in the repositories of
// one of the available package managers.
func SearchOnPackageManagers(packageName string) (results map[string]string, err error) {
	packageManagers := config.Vipers["config"].GetStringMapString("managers")

	for packageManager := range packageManagers {
		SearchOnPackageManager(packageManager, packageName)
	}
	return results, nil
}

// SearchOnPackageManager returns infos on packages found in the repositories of
// the given package manager.
func SearchOnPackageManager(packageManager string, packageName string) {
	fmt.Println("=======================")
	fmt.Printf("\n%s search %s", packageManager, packageName)
	termCmd := exec.Command(packageManager, "search", packageName)
	command.ExecuteCommand(termCmd)
}
