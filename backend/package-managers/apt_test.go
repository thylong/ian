package packagemanagers

import (
	"os/exec"
	"testing"
)

func TestAPTCmdWithArgs(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	apt := GetPackageManager("apt")

	cases := []struct {
		Method      func(string) error
		ExpectedErr error
	}{
		{apt.Install, nil},
		{apt.Uninstall, nil},
		{apt.UpdateOne, ErrAptMissingFeature},
		{apt.UpgradeOne, nil},
	}

	for _, tc := range cases {
		err := tc.Method("requests")
		if err != tc.ExpectedErr {
			t.Errorf("Expected nil error, got %#v", err)
		}
	}
}

func TestAPTCmdWithoutArgs(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	apt := GetPackageManager("apt")

	cases := []struct {
		Method      func() error
		ExpectedErr error
	}{
		{apt.Cleanup, nil},
		{apt.UpdateAll, nil},
		{apt.UpgradeAll, nil},
	}

	for _, tc := range cases {
		err := tc.Method()
		if err != tc.ExpectedErr {
			t.Errorf("Expected nil error, got %#v", err)
		}
	}
}

func TestAptGetExecPath(t *testing.T) {
	PackageManager := GetPackageManager("apt").(*AptPackageManager)

	if PackageManager.Path != PackageManager.GetExecPath() {
		t.Errorf("GetExecPath returned wrong Path: got %v want %v",
			PackageManager.GetExecPath(), PackageManager.Path)
	}
}
