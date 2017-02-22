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

package env

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/thylong/ian/backend/command"
)

// GetInfos returns env infos
func GetInfos() {
	IPCheckerURL := "http://httpbin.org/ip"

	resp, err := http.Get(IPCheckerURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error : %s", err.Error())
	}
	content, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var jsonContent map[string]string
	err = json.Unmarshal(content, &jsonContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error : %s", err.Error())
		return
	}

	command.ExecuteCommand(exec.Command("hostinfo"))
	fmt.Println("external_ip :", jsonContent["origin"])
	fmt.Print("uptime :")
	command.ExecuteCommand(exec.Command("uptime"))
}
