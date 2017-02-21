package packager

// PackageManager handles standard interactions with all Package Managers.
type PackageManager interface {
	Install(packageName string) (results string, err error)
	Uninstall(packageName string) (results string, err error)
	Cleanup() (results string, err error)
	Update() (results string, err error)
	Upgrade() (results string, err error)
	IsInstalled() bool
	IsOSPackageManager() bool
	GetExecPath() string
}
