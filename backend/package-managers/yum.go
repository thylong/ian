package packagemanagers

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/thylong/ian/backend/command"
)

// Yum immutable instance.
var Yum = YumPackageManager{Path: filepath.Clean("/usr/bin/yum"), Name: "yum"}

// YumPackageManager is the official Debian (and associated distributions) package manager.
// (more: https://wiki.debian.org/Yum)
type YumPackageManager struct {
	Path string
	Name string
}

// Install given Yum package.
func (yum *YumPackageManager) Install(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(yum.Path, "install", packageName)); err != nil {
		return fmt.Errorf("Cannot %s install %s: %s", yum.Name, packageName, err)
	}
	return err
}

// Uninstall given Yum package.
func (yum *YumPackageManager) Uninstall(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(yum.Path, "erase", packageName)); err != nil {
		return fmt.Errorf("Cannot %s erase %s: %s", yum.Name, packageName, err)
	}
	return err
}

// Cleanup all the local archives and previous versions.
func (yum *YumPackageManager) Cleanup() error {
	return command.ExecuteCommand(execCommand(yum.Path, "autoremove"))
}

// UpdateOne pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (yum *YumPackageManager) UpdateOne(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(yum.Path, "update", packageName)); err != nil {
		return fmt.Errorf("Cannot %s update %s: %s", yum.Name, packageName, err)
	}
	return err
}

// UpgradeOne Yum packages to the last known versions.
func (yum *YumPackageManager) UpgradeOne(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(yum.Path, "upgrade", packageName)); err != nil {
		return fmt.Errorf("Cannot %s upgrade %s: %s", yum.Name, packageName, err)
	}
	return err
}

// UpdateAll pulls last versions infos from realted repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (yum *YumPackageManager) UpdateAll() (err error) {
	if err = command.ExecuteCommand(execCommand(yum.Path, "update")); err != nil {
		return fmt.Errorf("Cannot %s update: %s", yum.Name, err)
	}
	return err
}

// UpgradeAll Yum packages to the last known versions.
func (yum *YumPackageManager) UpgradeAll() (err error) {
	if err = command.ExecuteCommand(execCommand(yum.Path, "upgrade")); err != nil {
		return fmt.Errorf("Cannot %s upgrade: %s", yum.Name, err)
	}
	return err
}

// IsInstalled returns true if Yum executable is found.
func (yum *YumPackageManager) IsInstalled() bool {
	if fileInfo, err := os.Stat(yum.Path); err != nil || fileInfo.Mode() == os.ModeSymlink {
		return false
	}
	return true
}

// IsOSPackageManager returns true for Mac OS.
func (yum *YumPackageManager) IsOSPackageManager() bool {
	return yum.IsInstalled() && runtime.GOOS == "linux"
}

// GetExecPath return immutable path to Yum executable.
func (yum *YumPackageManager) GetExecPath() string {
	return yum.Path
}

// GetName return the name of the package manager.
func (yum *YumPackageManager) GetName() string {
	return yum.Name
}

// Setup does nothing (yum comes by default in Linux distributions)
func (yum *YumPackageManager) Setup() (err error) {
	return nil
}
