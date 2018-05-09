package packagemanagers

import (
	"os/exec"
	"testing"
)

func TestYumCmdWithArgs(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	yum := GetPackageManager("yum")

	cases := []struct {
		Method      func(string) error
		ExpectedErr error
	}{
		{yum.Install, nil},
		{yum.Uninstall, nil},
		{yum.UpdateOne, nil},
		{yum.UpgradeOne, nil},
	}

	for _, tc := range cases {
		err := tc.Method("requests")
		if err != tc.ExpectedErr {
			t.Errorf("Expected nil error, got %#v", err)
		}
	}
}

func TestYumCmdWithoutArgs(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	yum := GetPackageManager("yum")

	cases := []struct {
		Method      func() error
		ExpectedErr error
	}{
		{yum.Cleanup, nil},
		{yum.UpdateAll, nil},
		{yum.UpgradeAll, nil},
	}

	for _, tc := range cases {
		err := tc.Method()
		if err != tc.ExpectedErr {
			t.Errorf("Expected nil error, got %#v", err)
		}
	}
}

func TestYumGetExecPath(t *testing.T) {
	PackageManager := GetPackageManager("yum").(*YumPackageManager)

	if PackageManager.Path != PackageManager.GetExecPath() {
		t.Errorf("GetExecPath returned wrong Path: got %v want %v",
			PackageManager.GetExecPath(), PackageManager.Path)
	}
}
