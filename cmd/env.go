// Copyright © 2016 THEOTIME LEVEQUE <theotime.leveque@gmail.com>
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

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/thylong/ian/backend/command"
)

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Get infos about the local environment",
	Long: `Get general or detailes informations about the current environment.
Currently implemented : System, Network, Security, current load.`,
}

var envInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get infos about the local environment",
	Long:  `Get general or detailes informations about the current environment. Currently implemented : System, Network, Security, current load.`,
	Run: func(cmd *cobra.Command, args []string) {
		IPCheckerURL := "http://httpbin.org/ip"

		resp, err := http.Get(IPCheckerURL)
		if err != nil {
			fmt.Println("Error : ", err.Error())
		}
		content, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		var jsonContent map[string]string
		err = json.Unmarshal(content, &jsonContent)
		if err != nil {
			fmt.Println("Error : ", err.Error())
			return
		}

		fmt.Println("====================")
		command.ExecuteCommand(exec.Command("hostinfo"))
		fmt.Println("====================")
		fmt.Println("external_ip :", jsonContent["origin"])
		fmt.Println("uptime :")
		command.ExecuteCommand(exec.Command("uptime"))
	},
}

var envUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the local environment",
	Long:  `Update the local environment with infos stored in the local config.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Updating env...")
		OSPackageManager.UpdateAll()
		OSPackageManager.UpgradeAll()
		OSPackageManager.Cleanup()
	},
}

func init() {
	RootCmd.AddCommand(envCmd)
	envCmd.AddCommand(envInfoCmd)
	envCmd.AddCommand(envUpdateCmd)
}
