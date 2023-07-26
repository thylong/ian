package packagemanagers

import (
	"os/exec"
	"testing"
)

func TestBrewCmdWithArgs(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	brew := GetPackageManager("brew")

	cases := []struct {
		Method      func(string) error
		ExpectedErr error
	}{
		{brew.Install, nil},
		{brew.Uninstall, nil},
		{brew.UpdateOne, nil},
		{brew.UpgradeOne, nil},
	}

	for _, tc := range cases {
		err := tc.Method("requests")
		if err != tc.ExpectedErr {
			t.Errorf("Expected nil error, got %#v", err)
		}
	}
}

func TestBrewCmdWithoutArgs(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	brew := GetPackageManager("brew")

	cases := []struct {
		Method      func() error
		ExpectedErr error
	}{
		{brew.Cleanup, nil},
		{brew.UpdateAll, nil},
		{brew.UpgradeAll, nil},
	}

	for _, tc := range cases {
		err := tc.Method()
		if err != tc.ExpectedErr {
			t.Errorf("Expected nil error, got %#v", err)
		}
	}
}

func TestBrewGetExecPath(t *testing.T) {
	PackageManager := GetPackageManager("brew").(*BrewPackageManager)

	if PackageManager.Path != PackageManager.GetExecPath() {
		t.Errorf("GetExecPath returned wrong Path: got %v want %v",
			PackageManager.GetExecPath(), PackageManager.Path)
	}
}
