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

package command

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

// ExecuteCommand a command and print output from stdout.
func ExecuteCommand(subCmd *exec.Cmd) (err error) {
	cmdReader, err := subCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Error creating StdoutPipe for Cmd: %v StdErr: %v", err, os.Stderr)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
	}()

	err = subCmd.Start()
	if err != nil {
		return fmt.Errorf("Error starting Cmd: %v StdErr: %v", err, os.Stderr)
	}

	err = subCmd.Wait()
	if err != nil {
		return fmt.Errorf("Error waiting for Cmd: %v StdErr: %v", err, os.Stderr)
	}

	return nil
}
