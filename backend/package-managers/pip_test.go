// Copyright © 2016 Theotime LEVEQUE theotime@protonmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	PackageManager := GetPackageManager("pip").(PipPackageManager)

	if PackageManager.Path != PackageManager.GetExecPath() {
		t.Errorf("GetExecPath returned wrong Path: got %v want %v",
			PackageManager.GetExecPath(), PackageManager.Path)
	}
}
