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

package env

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"

	"github.com/fatih/color"
	"github.com/thylong/ian/backend/command"
	pm "github.com/thylong/ian/backend/package-managers"
)

// SetupDotFiles ask and retrieve a dotfiles repository.
func SetupDotFiles(dotfilesRepository string, dotfilesDirPath string) {
	usr, _ := user.Current()
	if _, err := os.Stat(usr.HomeDir + "/.dotfiles"); err != nil && dotfilesRepository != "" {
		termCmd := exec.Command("git", "clone", "-v", "https://github.com/"+dotfilesRepository+".git", dotfilesDirPath)
		command.ExecuteInteractiveCommand(termCmd)

		re := regexp.MustCompile(".git$")

		files, _ := ioutil.ReadDir(filepath.Join(usr.HomeDir, "/.dotfiles"))
		for _, f := range files {
			if re.MatchString(f.Name()) {
				continue
			}

			if _, err := os.Stat(filepath.Join(usr.HomeDir, f.Name())); err != nil {
				if err := os.Symlink(
					filepath.Join(usr.HomeDir, ".dotfiles", f.Name()),
					filepath.Join(usr.HomeDir, f.Name()),
				); err != nil {
					fmt.Fprintf(os.Stderr, "%v %s.", color.RedString("Error:"), err)
				}
			}
		}
	} else {
		fmt.Println("Skipping dotfiles configuration.")
	}
}

// SetupPackages installs listed CLI packages.
func SetupPackages(PackageManager pm.PackageManager, packages []string) {
	fmt.Println("Installing packages...")

	if len(packages) == 0 {
		fmt.Println("No packages to install")
		return
	}

	for _, packageToInstall := range packages {
		if err := PackageManager.Install(packageToInstall); err != nil {
			fmt.Fprintf(os.Stderr, "%v %s.", color.RedString("Error:"), err)
		}
	}
}
