package packagemanagers

import (
	"os/exec"
	"testing"
)

func TestNPMCmdWithArgs(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	npm := GetPackageManager("npm")

	cases := []struct {
		Method      func(string) error
		ExpectedErr error
	}{
		{npm.Install, nil},
		{npm.Uninstall, nil},
		{npm.UpdateOne, nil},
		{npm.UpgradeOne, nil},
	}

	for _, tc := range cases {
		err := tc.Method("requests")
		if err != tc.ExpectedErr {
			t.Errorf("Expected nil error, got %#v", err)
		}
	}
}

func TestNPMCmdWithoutArgs(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	npm := GetPackageManager("npm")

	cases := []struct {
		Method      func() error
		ExpectedErr error
	}{
		{npm.Cleanup, nil},
		{npm.UpdateAll, ErrNPMMissingFeature},
		{npm.UpgradeAll, nil},
	}

	for _, tc := range cases {
		err := tc.Method()
		if err != tc.ExpectedErr {
			t.Errorf("Expected nil error, got %#v", err)
		}
	}
}

func TestNPMGetExecPath(t *testing.T) {
	PackageManager := GetPackageManager("npm").(*NpmPackageManager)

	if PackageManager.Path != PackageManager.GetExecPath() {
		t.Errorf("GetExecPath returned wrong Path: got %v want %v",
			PackageManager.GetExecPath(), PackageManager.Path)
	}
}
