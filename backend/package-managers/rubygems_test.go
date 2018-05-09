package packagemanagers

import (
	"os/exec"
	"testing"
)

func TestRubyGemsCmdWithArgs(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	rubygems := GetPackageManager("rubygems")

	cases := []struct {
		Method      func(string) error
		ExpectedErr error
	}{
		{rubygems.Install, nil},
		{rubygems.Uninstall, nil},
		{rubygems.UpdateOne, ErrRubyGemsMissingFeature},
		{rubygems.UpgradeOne, nil},
	}

	for _, tc := range cases {
		err := tc.Method("requests")
		if err != tc.ExpectedErr {
			t.Errorf("Expected nil error, got %#v", err)
		}
	}
}

func TestRubyGemsCmdWithoutArgs(t *testing.T) {
	execCommand = mockExecCommand
	defer func() { execCommand = exec.Command }()

	rubygems := GetPackageManager("rubygems")

	cases := []struct {
		Method      func() error
		ExpectedErr error
	}{
		{rubygems.Cleanup, nil},
		{rubygems.UpdateAll, ErrRubyGemsMissingFeature},
		{rubygems.UpgradeAll, nil},
	}

	for _, tc := range cases {
		err := tc.Method()
		if err != tc.ExpectedErr {
			t.Errorf("Expected nil error, got %#v", err)
		}
	}
}

func TestRubyGemsGetExecPath(t *testing.T) {
	PackageManager := GetPackageManager("rubygems").(*RubyGemsPackageManager)

	if PackageManager.Path != PackageManager.GetExecPath() {
		t.Errorf("GetExecPath returned wrong Path: got %v want %v",
			PackageManager.GetExecPath(), PackageManager.Path)
	}
}
