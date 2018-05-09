package packagemanagers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/thylong/ian/backend/command"
)

// Npm immutable instance.
var Npm = NpmPackageManager{Path: filepath.Clean("/usr/local/bin/npm"), Name: "npm"}

// ErrNPMMissingFeature is returned when triggering an unsupported feature.
var ErrNPMMissingFeature = errors.New("npm is not designed to support this feature")

// NpmPackageManager is a (widely used) unofficial Mac OS package manager.
// (more: https://npm.sh/)
type NpmPackageManager struct {
	Path string
	Name string
}

// Install given Npm package.
func (npm *NpmPackageManager) Install(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(npm.Path, "install", "-g", packageName)); err != nil {
		return fmt.Errorf("Cannot %s install %s: %s", npm.Name, packageName, err)
	}
	return err
}

// Uninstall given Npm package.
func (npm *NpmPackageManager) Uninstall(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(npm.Path, "uninstall", "-g", packageName)); err != nil {
		return fmt.Errorf("Cannot %s install %s: %s", npm.Name, packageName, err)
	}
	return err
}

// Cleanup the npm cache.
func (npm *NpmPackageManager) Cleanup() error {
	return command.ExecuteCommand(execCommand(npm.Path, "cache", "clean"))
}

// UpdateOne pulls last versions infos from related repositories.
// This is not performing any updates and should be coupled
// with upgradeAll command.
func (npm *NpmPackageManager) UpdateOne(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(npm.Path, "update", packageName)); err != nil {
		return fmt.Errorf("Cannot %s update %s: %s", npm.Name, packageName, err)
	}
	return err
}

// UpgradeOne Npm packages to the last known versions.
func (npm *NpmPackageManager) UpgradeOne(packageName string) (err error) {
	if err = command.ExecuteCommand(execCommand(npm.Path, "upgrade", packageName)); err != nil {
		return fmt.Errorf("Cannot %s upgrade %s: %s", npm.Name, packageName, err)
	}
	return err
}

// UpdateAll does nothing (out of making NPM satisfying PackageManager interface).
func (npm *NpmPackageManager) UpdateAll() error {
	return ErrNPMMissingFeature
}

// UpgradeAll Npm packages to the last known versions.
func (npm *NpmPackageManager) UpgradeAll() (err error) {
	if err = command.ExecuteCommand(execCommand(npm.Path, "update", "-g")); err != nil {
		return fmt.Errorf("Cannot %s update: %s", npm.Name, err)
	}
	return err
}

// IsInstalled returns true if Npm executable is found.
func (npm *NpmPackageManager) IsInstalled() bool {
	if _, err := os.Stat(npm.Path); err != nil {
		return false
	}
	return true
}

// IsOSPackageManager returns false because npm is never the main OS Package Manager.
func (npm *NpmPackageManager) IsOSPackageManager() bool {
	return false
}

// GetExecPath return immutable path to Npm executable.
func (npm *NpmPackageManager) GetExecPath() string {
	return npm.Path
}

// GetName return the name of the package manager.
func (npm *NpmPackageManager) GetName() string {
	return npm.Name
}

// Setup installs Cask
func (npm *NpmPackageManager) Setup() (err error) {
	return nil
}
