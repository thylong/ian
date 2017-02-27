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
	"net/http"
	"os/exec"

	"github.com/thylong/ian/backend/command"
	"github.com/thylong/ian/backend/config"
)

// PackageManager handles standard interactions with all Package Managers.
type PackageManager interface {
	Install(packageName string) (err error)
	Uninstall(packageName string) (err error)
	Cleanup() (err error)
	UpdateOne(string) (err error)
	UpgradeOne(string) (err error)
	UpdateAll() (err error)
	UpgradeAll() (err error)
	IsInstalled() bool
	IsOSPackageManager() bool
	GetExecPath() string
	GetName() string
	Setup() (err error)
}

var supportedPackageManagers = make(map[string]PackageManager)

func init() {
	supportedPackageManagers["brew"] = Brew
	supportedPackageManagers["pip"] = Pip
	supportedPackageManagers["npm"] = Npm
	supportedPackageManagers["apt"] = Apt
	supportedPackageManagers["yum"] = Yum
	supportedPackageManagers["rubygems"] = RubyGems
}

// GetOSPackageManager returns the main Package Manager of the current OS.
// As only MacOS is supported for now, it returns a Brew instance.
func GetOSPackageManager() PackageManager {
	return Brew
}

// GetPackageManager returns the corresponding PackageManager otherwise default
// to OS package manager.
func GetPackageManager(PackageManagerFlag string) PackageManager {
	packageManager, ok := supportedPackageManagers[PackageManagerFlag]

	if !ok {
		return packageManager
	}
	return GetOSPackageManager()
}

// IsAvailableOnPackageManagers returns true if found in the repositories of
// one of the available package managers.
func IsAvailableOnPackageManagers(packageName string) (map[string]bool, error) {
	packageManagers := config.Vipers["config"].GetStringMap("managers")
	results := make(map[string]bool)

	for packageManager, packageParams := range packageManagers {
		baseURL := packageParams.(map[interface{}]interface{})["base_url"].(string)
		results[packageManager] = isAvailableOnPackageManager(packageManager, baseURL, packageName)
	}
	return results, nil
}

// IsAvailableOnPackageManagers returns true if found in the repositories of
// the given package manager.
func isAvailableOnPackageManager(packageManager string, baseURL string, packageName string) bool {
	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if packageManager == "composer" && req.URL.String() != baseURL+packageName+"/" {
			return errors.New("Fail on redirect...")
		}
		return nil
	}

	resp, err := client.Head(baseURL + packageName)
	if err != nil || resp.StatusCode != 200 {
		fmt.Printf("%s is not reachable.", packageManager)
		return false
	}
	return true
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
	fmt.Print("=======================")
	fmt.Printf("%s search %s", packageManager, packageName)
	fmt.Print("=======================")
	termCmd := exec.Command(packageManager, "search", packageName)
	command.ExecuteCommand(termCmd)
}
