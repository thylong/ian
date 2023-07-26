package packagemanagers

import (
	"os/exec"
	"testing"
)

func TestPIPCmdWithArgs(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	pip := GetPackageManager("pip")

	cases := []struct {
		Method      func(string) error
		ExpectedErr error
	}{
		{pip.Install, nil},
		{pip.Uninstall, nil},
		{pip.UpdateOne, ErrPipMissingFeature},
		{pip.UpgradeOne, nil},
	}

	for _, tc := range cases {
		err := tc.Method("requests")
		if err != tc.ExpectedErr {
			t.Errorf("Expected nil error, got %#v", err)
		}
	}
}

func TestPIPCmdWithoutArgs(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	pip := GetPackageManager("pip")

	cases := []struct {
		Method      func() error
		ExpectedErr error
	}{
		{pip.Cleanup, ErrPipMissingFeature},
		{pip.UpdateAll, ErrPipMissingFeature},
		{pip.UpgradeAll, ErrPipMissingFeature},
	}

	for _, tc := range cases {
		err := tc.Method()
		if err != tc.ExpectedErr {
			t.Errorf("Expected nil error, got %#v", err)
		}
	}
}

func TestPIPGetExecPath(t *testing.T) {
	PackageManager := GetPackageManager("pip").(*PipPackageManager)

	if PackageManager.Path != PackageManager.GetExecPath() {
		t.Errorf("GetExecPath returned wrong Path: got %v want %v",
			PackageManager.GetExecPath(), PackageManager.Path)
	}
}
