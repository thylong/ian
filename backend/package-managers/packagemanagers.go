package packagemanagers

import (
	"errors"
	"os/exec"
)

// PackageManager handles standard interactions with all Package Managers.
type PackageManager interface {
	Install(packageName string) error
	Uninstall(packageName string) error
	Cleanup() error
	UpdateOne(string) error
	UpgradeOne(string) error
	UpdateAll() error
	UpgradeAll() error
	IsInstalled() bool
	IsOSPackageManager() bool
	GetExecPath() string
	GetName() string
	Setup() error
}

// SupportedPackageManagers contains all the currently supported package managers.
var SupportedPackageManagers = make(map[string]PackageManager)

var execCommand = exec.Command

func init() {
	SupportedPackageManagers["brew"] = &Brew
	SupportedPackageManagers["cask"] = &Cask
	SupportedPackageManagers["pip"] = &Pip
	SupportedPackageManagers["npm"] = &Npm
	SupportedPackageManagers["apt"] = &Apt
	SupportedPackageManagers["yum"] = &Yum
	SupportedPackageManagers["rubygems"] = &RubyGems
	SupportedPackageManagers["apm"] = &Apm
}

// GetOSPackageManager returns the main Package Manager of the current OS.
// As only MacOS is supported for now, it returns a Brew instance.
func GetOSPackageManager() (PackageManager, error) {
	for name, packageManager := range SupportedPackageManagers {
		if name != "cask" && packageManager.IsOSPackageManager() {
			return packageManager, nil
		}
	}
	return &Brew, errors.New("No OS Package Manager found")
}

// GetPackageManager returns the corresponding PackageManager otherwise default
// to OS package manager.
func GetPackageManager(PackageManagerFlag string) PackageManager {
	packageManager, ok := SupportedPackageManagers[PackageManagerFlag]
	if ok {
		return packageManager
	}

	OSPackageManager, _ := GetOSPackageManager()
	return OSPackageManager
}

// UpdateAllPackageManagers updates all packages managers.
func UpdateAllPackageManagers() {
	for _, packageManager := range SupportedPackageManagers {
		if packageManager.IsInstalled() {
			packageManager.UpdateAll()
		}
	}
}

// UpgradeAllPackageManagers upgrades all packages from package managers.
func UpgradeAllPackageManagers() {
	for _, packageManager := range SupportedPackageManagers {
		if packageManager.IsInstalled() {
			packageManager.UpdateAll()
		}
	}
}

// IsSupportedPackageManager returns true if the PackageManager is supported by
// Ian else returns false.
func IsSupportedPackageManager(packageManager string) bool {
	if _, ok := SupportedPackageManagers[packageManager]; !ok {
		return false
	}
	return true
}
