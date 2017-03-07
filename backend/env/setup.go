package env

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"regexp"

	"github.com/thylong/ian/backend/command"
	pm "github.com/thylong/ian/backend/package-managers"
)

// SetupPackages installs listed CLI packages.
func SetupPackages(PackageManager pm.PackageManager, packages []string) {
	fmt.Println("Installing packages...")

	if len(packages) == 0 {
		fmt.Println("No packages to install")
		return
	}

	for _, packageToInstall := range packages {
		PackageManager.Install(packageToInstall)
	}
}

// SetupDotFiles ask for a Github nickname and retrieve the dotfiles repo
// (the repository has to be public).
func SetupDotFiles(nickname string, dotfilesDirPath string) {
	usr, _ := user.Current()
	if _, err := os.Stat(usr.HomeDir + "/.dotfiles"); err != nil && nickname != "" {
		termCmd := exec.Command("git", "clone", "-v", "git@github.com:"+nickname+"/dotfiles.git", dotfilesDirPath)
		command.ExecuteCommand(termCmd)

		re := regexp.MustCompile(".git$")

		files, _ := ioutil.ReadDir(usr.HomeDir + "/.dotfiles")
		for _, f := range files {
			if re.MatchString(f.Name()) {
				continue
			}

			if _, err := os.Stat(usr.HomeDir + "/" + f.Name()); err != nil {
				err := os.Symlink(usr.HomeDir+"/.dotfiles/"+f.Name(), usr.HomeDir+"/"+f.Name())
				if err != nil {
					fmt.Fprint(os.Stderr, err)
				}
			}
		}
	} else {
		fmt.Println("Skipping dotfiles configuration.")
	}
}
