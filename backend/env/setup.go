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

// SetupCLIPackages installs listed CLI packages.
func SetupCLIPackages(OSPackageManager pm.PackageManager, CLIPackages []string) {
	fmt.Println("Installing CLI packages...")

	if len(CLIPackages) == 0 {
		fmt.Println("No brew packages to install")
		return
	}

	for _, CLIPackage := range CLIPackages {
		OSPackageManager.Install(CLIPackage)
	}
}

// SetupGUIPackages installs listed GUI packages.
func SetupGUIPackages(OSPackageManager pm.PackageManager, GUIPackages []string) {
	fmt.Println("Installing GUI packages...")

	if len(GUIPackages) == 0 {
		fmt.Println("No GUI packages to install")
		return
	}

	for _, GUIPackage := range GUIPackages {
		if OSPackageManager.GetName() == "brew" {
			command.ExecuteCommand(exec.Command(OSPackageManager.GetExecPath(), "cask", "install", GUIPackage))
		} else {
			OSPackageManager.Install(GUIPackage)
		}
	}
}

// SetupDotFiles ask for a Github nickname and retrieve the dotfiles repo
// (the repository has to be public).
func SetupDotFiles(nickname string, dotfilesDirPath string) {
	if nickname != "" {
		termCmd := exec.Command("git", "clone", "-v", "git@github.com:"+nickname+"/dotfiles.git", dotfilesDirPath)
		command.ExecuteCommand(termCmd)

		re := regexp.MustCompile(".git$")

		usr, _ := user.Current()
		files, _ := ioutil.ReadDir(usr.HomeDir + "/.dotfiles")
		for _, f := range files {
			if re.MatchString(f.Name()) {
				continue
			}
			if err := os.Symlink(usr.HomeDir+"/.dotfiles/"+f.Name(), usr.HomeDir+"/"+f.Name()); err != nil {
				fmt.Fprint(os.Stderr, err)
			}
		}
	} else {
		fmt.Print("Skipping dotfiles configuration.")
	}
}
