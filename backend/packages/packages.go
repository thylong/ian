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

package packages

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thylong/ian/backend/command"
	"github.com/thylong/ian/backend/config"
	yaml "gopkg.in/yaml.v2"
)

// IsAvailableOnPackageManagers returns true if found in the repositories of
// one of the available package managers.
func IsAvailableOnPackageManagers(packageName string) (map[string]bool, error) {
	packageManagers := viper.GetStringMap("managers")
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
		log.Debug(packageManager + " is not reachable.")
		return false
	}
	return true
}

// SearchOnPackageManagers returns infos on packages found in the repositories of
// one of the available package managers.
func SearchOnPackageManagers(packageName string) (results map[string]string, err error) {
	packageManagers := viper.GetStringMapString("managers")

	for packageManager := range packageManagers {
		SearchOnPackageManager(packageManager, packageName)
	}
	return results, nil
}

// SearchOnPackageManager returns infos on packages found in the repositories of
// the given package manager.
func SearchOnPackageManager(packageManager string, packageName string) {
	fmt.Println("=======================")
	fmt.Println(packageManager, "search", packageName)
	fmt.Println("=======================")
	termCmd := exec.Command(packageManager, "search", packageName)
	command.ExecuteCommand(termCmd)
}

// WritePackageEntry in the local config.yml
func WritePackageEntry(selectedPackageManager string, arg string) error {
	config.Config.Packages[arg] = map[string]string{
		"cmd":         arg,
		"description": arg,
		"type":        selectedPackageManager,
	}
	ymlContent, _ := yaml.Marshal(config.Config)
	err := ioutil.WriteFile(config.ConfigFullPath, ymlContent, 0666)
	if err != nil {
		return errors.New("Unable to edit config file.")
	}
	return nil
}

// UnwritePackageEntry in the local config.yml
func UnwritePackageEntry(selectedPackageManager string, arg string) error {
	delete(config.Config.Packages, arg)
	ymlContent, _ := yaml.Marshal(config.Config)
	err := ioutil.WriteFile(config.ConfigFullPath, ymlContent, 0666)
	if err != nil {
		return errors.New("Unable to edit config file.")
	}
	return nil
}
