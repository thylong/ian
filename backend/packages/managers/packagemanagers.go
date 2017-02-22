package packagemanagers

// PackageManager handles standard interactions with all Package Managers.
type PackageManager interface {
	Install(packageName string) (err error)
	Uninstall(packageName string) (err error)
	Cleanup() (err error)
	UpdateAll() (err error)
	UpgradeAll() (err error)
	IsInstalled() bool
	IsOSPackageManager() bool
	GetExecPath() string
	Setup() (err error)
}

// GetOSPackageManager returns the main Package Manager of the current OS.
// As only MacOS is supported for now, it returns a Brew instance.
func GetOSPackageManager() PackageManager {
	return Brew
}
