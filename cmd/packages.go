// Copyright © 2016 Theotime Leveque <theotime.leveque@gmail.com>
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

package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	yaml "gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// packagesCmd represents the packages command
var packagesCmd = &cobra.Command{
	Use:   "packages",
	Short: "packages allows you to manage ian extensions",
	Long: `packages allows you to install, uninstall subcommands for Ian.

    An example would be baily CLI (a nice bot powered by @samueldelesque & Dailymotion)`,
}

var selectedPackageManager string

func init() {
	installPackagesCmd.Flags().StringVarP(&selectedPackageManager, "package-manager", "p", "", "Package manager to use to install.")
	uninstallPackagesCmd.Flags().StringVarP(&selectedPackageManager, "package-manager", "p", "", "Package manager to use to install.")
	RootCmd.AddCommand(packagesCmd)
	packagesCmd.AddCommand(listPackagesCmd)
	packagesCmd.AddCommand(installPackagesCmd)
	packagesCmd.AddCommand(uninstallPackagesCmd)
	packagesCmd.AddCommand(searchPackagesCmd)
}

var listPackagesCmd = &cobra.Command{
	Use:   "list",
	Short: "list installed ian extensions",
	Long: `list installed ian extensions,

    This won't list npm, pip, gem, composer or other kind of packages.`,
	Run: func(cmd *cobra.Command, args []string) {
		packagesUsages := `Package Commands:`
		for packageName, packageMeta := range Config.Packages {
			packagesUsages += "\n" + `  ` + packageName + ` ` + packageMeta["description"] + ` type:` + packageMeta["type"]
		}
		fmt.Println(packagesUsages)
	},
}

var installPackagesCmd = &cobra.Command{
	Use:   "install",
	Short: "install a new extension",
	Long: `install new subcommand(s) for Ian.

    An example would be baily CLI (a nice bot powered by @samueldelesque & Dailymotion)`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			results, err := isAvailableOnPackageManagers(arg)
			if err != nil {
				log.Error(err.Error())
				return
			}

			if selectedPackageManager == "" {
				for packageManager, available := range results {
					if available {
						fmt.Println("Package found on :", packageManager)
					}
				}
				fmt.Println("\nUse -p option to install,\nian packages --help to print usage.")
			} else {
				cmdParams := []string{}
				installParams := strings.Split(Config.Managers[selectedPackageManager]["install_cmd"], " ")
				cmdParams = append(installParams, arg)

				termCmd := exec.Command(selectedPackageManager, cmdParams...)
				printFromCmdStds(termCmd)

				err := writePackageEntry(selectedPackageManager, arg)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	},
}

var searchPackagesCmd = &cobra.Command{
	Use:   "search",
	Short: "search a new extension",
	Long: `search subcommand(s) for Ian.

    An example would be baily CLI (a nice bot powered by @samueldelesque & Dailymotion)`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			results, err := searchOnPackageManagers(arg)
			if err != nil {
				log.Error(err.Error())
				return
			}
			fmt.Println(results)
		}
	},
}

var uninstallPackagesCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "uninstall an extension",
	Long:  `uninstall an extension.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			results, err := isAvailableOnPackageManagers(arg)
			if err != nil {
				log.Error(err.Error())
				return
			}

			if selectedPackageManager == "" {
				for packageManager, available := range results {
					if available {
						fmt.Println("Package found on :", packageManager)
					}
				}
				fmt.Println("\nUse -p option to uninstall,\nian packages --help to print usage.")
			} else {
				cmdParams := []string{}
				installParams := strings.Split(Config.Managers[selectedPackageManager]["uninstall_cmd"], " ")
				cmdParams = append(installParams, arg)
				fmt.Println(cmdParams)

				termCmd := exec.Command(selectedPackageManager, cmdParams...)
				printFromCmdStds(termCmd)

				err := unwritePackageEntry(selectedPackageManager, arg)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	},
}

func isAvailableOnPackageManagers(packageName string) (map[string]bool, error) {
	packageManagers := viper.GetStringMap("managers")
	results := make(map[string]bool)

	for packageManager, packageParams := range packageManagers {
		baseURL := packageParams.(map[interface{}]interface{})["base_url"].(string)
		results[packageManager] = isAvailableOnPackageManager(packageManager, baseURL, packageName)
	}
	return results, nil
}

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

func searchOnPackageManagers(packageName string) (results map[string]string, err error) {
	packageManagers := viper.GetStringMapString("managers")

	for packageManager := range packageManagers {
		searchOnPackageManager(packageManager, packageName)
	}
	return results, nil
}

func searchOnPackageManager(packageManager string, packageName string) {
	fmt.Println("=======================")
	fmt.Println(packageManager, "search", packageName)
	fmt.Println("=======================")
	termCmd := exec.Command(packageManager, "search", packageName)
	printFromCmdStds(termCmd)
}

func writePackageEntry(selectedPackageManager string, arg string) error {
	Config.Packages[arg] = map[string]string{
		"cmd":         arg,
		"description": arg,
		"type":        selectedPackageManager,
	}
	ymlContent, _ := yaml.Marshal(Config)
	err := ioutil.WriteFile(ConfigFullPath, ymlContent, 0666)
	if err != nil {
		return errors.New("Unable to edit config file.")
	}
	return nil
}

func unwritePackageEntry(selectedPackageManager string, arg string) error {
	delete(Config.Packages, arg)
	ymlContent, _ := yaml.Marshal(Config)
	err := ioutil.WriteFile(ConfigFullPath, ymlContent, 0666)
	if err != nil {
		return errors.New("Unable to edit config file.")
	}
	return nil
}
