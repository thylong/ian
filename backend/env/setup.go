package env

import (
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"

	"github.com/thylong/ian/backend/command"
	"github.com/thylong/ian/backend/config"
	"github.com/thylong/ian/backend/log"
	pm "github.com/thylong/ian/backend/package-managers"
)

// Setup installs Ian and configuration Ian's environment.
func Setup(OSPackageManager pm.PackageManager) {
	if _, err := os.Stat(OSPackageManager.GetExecPath()); err != nil {
		log.Infoln("Installing OS package manager...")
		if err = OSPackageManager.Setup(); err != nil {
			log.Errorln("Missing OS package manager !")
			return
		}
	}

	SetupDotFiles(
		config.Vipers["config"].GetStringMapString("dotfiles")["repository"],
		config.DotfilesDirPath,
	)

	// Refresh the configuration in case the imported dotfiels contains ian configuration
	config.RefreshVipers()

	log.Warningln("You don't have any packages to be installed in your current ian configuration.")
	if _, ok := config.Vipers["env"]; !ok && config.GetBoolUserInput("Would you like to use a preset? (Y/n)") {
		in := config.GetUserInput(`Which preset would you like to use:
1) Software engineer (generalist preset)
2) Backend developer
3) Frontend developer
4) Ops
Enter your choice`)
		config.CreateEnvFileWithPreset(in)
	}

	packageManagers := config.Vipers["env"].AllKeys()
	for _, packageManager := range packageManagers {
		packages := config.Vipers["env"].GetStringSlice(packageManager)
		InstallPackages(pm.GetPackageManager(packageManager), packages)
	}
}

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
					log.Errorln(err)
				}
			}
		}
	} else {
		log.Infoln("Skipping dotfiles configuration.")
	}
}

// InstallPackages installs listed CLI packages.
func InstallPackages(PackageManager pm.PackageManager, packages []string) {
	if len(packages) == 0 {
		return
	}
	log.Infof("Installing %s packages...", PackageManager.GetName())

	for _, packageToInstall := range packages {
		if err := PackageManager.Install(packageToInstall); err != nil {
			log.Errorln(err)
		}
	}
}
