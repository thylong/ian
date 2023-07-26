package packagemanagers

import (
	"os/exec"
	"testing"
)

func TestCaskCmdWithArgs(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	cask := GetPackageManager("cask")

	cases := []struct {
		Method      func(string) error
		ExpectedErr error
	}{
		{cask.Install, nil},
		{cask.Uninstall, nil},
		{cask.UpdateOne, nil},
		{cask.UpgradeOne, nil},
	}

	for _, tc := range cases {
		err := tc.Method("requests")
		if err != tc.ExpectedErr {
			t.Errorf("Expected nil error, got %#v", err)
		}
	}
}

func TestCaskCmdWithoutArgs(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	cask := GetPackageManager("cask")

	cases := []struct {
		Method      func() error
		ExpectedErr error
	}{
		{cask.Cleanup, nil},
		{cask.UpdateAll, nil},
		{cask.UpgradeAll, nil},
	}

	for _, tc := range cases {
		err := tc.Method()
		if err != tc.ExpectedErr {
			t.Errorf("Expected nil error, got %#v", err)
		}
	}
}

func TestCaskGetExecPath(t *testing.T) {
	PackageManager := GetPackageManager("cask").(*CaskPackageManager)

	if PackageManager.Path != PackageManager.GetExecPath() {
		t.Errorf("GetExecPath returned wrong Path: got %v want %v",
			PackageManager.GetExecPath(), PackageManager.Path)
	}
}
