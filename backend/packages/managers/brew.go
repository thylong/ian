package packagemanagers

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/thylong/ian/backend/command"
)

// Brew immutable instance.
var Brew = BrewPackageManager{Path: "/usr/local/bin/brew"}

// BrewPackageManager is a (widely used) unofficial Mac OS package manager.
// (more: https://brew.sh/)
type BrewPackageManager struct {
	Path string
}

// Install given Brew package.
func (b BrewPackageManager) Install(packageName string) (err error) {
	err = command.ExecuteCommand(exec.Command(b.Path, "install", packageName))
	return err
}

// Uninstall given Brew package.
func (b BrewPackageManager) Uninstall(packageName string) (err error) {
	err = command.ExecuteCommand(exec.Command(b.Path, "uninstall", packageName))
	return err
}

// Cleanup all the local archives and previous versions.
func (b BrewPackageManager) Cleanup() (err error) {
	// Cleanup brew
	err = command.ExecuteCommand(exec.Command(b.Path, "cleanup"))
	// Cleanup cask
	command.ExecuteCommand(exec.Command("/usr/local/bin/brew", "cask", "cleanup"))
	return err
}

// UpdateAll pulls last versions infos from realted repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (b BrewPackageManager) UpdateAll() (err error) {
	err = command.ExecuteCommand(exec.Command(b.Path, "update"))
	return err
}

// UpgradeAll Brew packages to the last known versions.
func (b BrewPackageManager) UpgradeAll() (err error) {
	err = command.ExecuteCommand(exec.Command(b.Path, "upgrade"))
	return err
}

// IsInstalled returns true if Brew executable is found.
func (b BrewPackageManager) IsInstalled() bool {
	if _, err := os.Stat(b.Path); err != nil {
		return false
	}
	return true
}

// IsOSPackageManager returns true for Mac OS.
func (b BrewPackageManager) IsOSPackageManager() bool {
	return b.IsInstalled()
}

// GetExecPath return immutable path to Brew executable.
func (b BrewPackageManager) GetExecPath() string {
	return b.Path
}

// Setup installs Cask
func (b BrewPackageManager) Setup() (err error) {
	fmt.Println("Installing cask...")
	if _, err := os.Stat("/usr/local/bin/cask"); err != nil {
		err = command.ExecuteCommand(exec.Command("brew", "tap", "caskroom/cask"))
		return err
	}
	fmt.Println("cask already installed, skipping...")
	return nil
}
