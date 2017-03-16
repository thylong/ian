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
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func TestGetOSPackageManager(t *testing.T) {
	OS := runtime.GOOS
	if OS == "windows" {
		t.Skip("Windows is not supported yet.")
	}

	OSPackageManager, _ := GetOSPackageManager()
	OSPackageManagerName := OSPackageManager.GetName()

	if OS == "darwin" && OSPackageManagerName != "brew" {
		t.Errorf("OS returned wrong OS package manager: got %v want %v",
			OSPackageManagerName, "brew")
	}

	if OS == "linux" {
		if fileContent, err := ioutil.ReadFile("/etc/issue"); err == nil {
			if strings.Contains("Ubuntu", string(fileContent)) && OSPackageManagerName != "apt" {
				t.Errorf("OS returned wrong OS package manager: got %v want %v",
					OSPackageManagerName, "apt")
			}
			if strings.Contains("CentOS", string(fileContent)) && OSPackageManagerName != "yum" {
				t.Errorf("OS returned wrong OS package manager: got %v want %v",
					OSPackageManagerName, "yum")
			}
			return
		}
		t.Skip("Linux distribution not supported yet.")
	}
	t.Skip("Current OS is not supported yet.")
}

func TestGetPackageManager(t *testing.T) {
	for PackageManagerName := range SupportedPackageManagers {
		if pm := GetPackageManager(PackageManagerName).GetName(); pm != PackageManagerName {
			t.Errorf("GetPackageManager returned wrong Package manager: got %v want %v",
				pm, PackageManagerName)
		}
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// fmt.Println(os.Args[3:])
	os.Exit(0)
}

func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}
