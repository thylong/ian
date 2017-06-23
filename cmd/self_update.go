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

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/mitchellh/ioprogress"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var httpGet = http.Get
var version string
var ianLastRelease = "https://api.github.com/repos/thylong/ian/releases/latest"

func init() {
	RootCmd.AddCommand(selfUpdateCmd)
}

var selfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Update ian to the last version",
	Long:  `Update ian to the last version.`,
	Run: func(cmd *cobra.Command, args []string) {
		localVersion := viper.GetString("VERSION")

		resp, err := httpGet(ianLastRelease)
		if err != nil {
			fmt.Printf("%v Cannot retrieve ian's last release infos\n", color.RedString("Error:"))
			return
		}
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%v Cannot retrieve ian's last release infos\n", color.RedString("Error:"))
			return
		}
		defer resp.Body.Close()

		var lastReleaseContent map[string]interface{}
		err = json.Unmarshal(content, &lastReleaseContent)
		if err != nil {
			fmt.Printf("%v Cannot retrieve ian's last version\n", color.RedString("Error:"))
			return
		}
		remoteVersion := lastReleaseContent["tag_name"]

		if localVersion != remoteVersion {
			fmt.Printf("ian version outdated... \n")
			downloadURL := lastReleaseContent["assets"].([]interface{})[0].(map[string]interface{})["browser_download_url"]
			resp, err = httpGet(downloadURL.(string))
			if err != nil {
				fmt.Printf("%v Cannot download ian's last version\n", color.RedString("Error:"))
				return
			}
			defer resp.Body.Close()

			progressR := &ioprogress.Reader{
				Reader:       resp.Body,
				Size:         resp.ContentLength,
				DrawInterval: 500 * time.Millisecond,
				DrawFunc: ioprogress.DrawTerminalf(os.Stdout, func(progress, total int64) string {
					bar := ioprogress.DrawTextFormatBar(40)
					return fmt.Sprintf("%s %20s", bar(progress, total), ioprogress.DrawTextFormatBytes(progress, total))
				}),
			}

			data, err := ioutil.ReadAll(progressR)
			if err != nil {
				fmt.Printf("%v ian's last version seems broken\n", color.RedString("Error:"))
				return
			}
			dest, err := os.Executable()
			if err != nil {
				fmt.Printf("%v ian's last version seems not executable\n", color.RedString("Error:"))
				return
			}

			// Move the old version to a backup path that we can recover from
			// in case the upgrade fails
			destBackup := dest + ".bak"
			if _, err := os.Stat(dest); err == nil {
				os.Rename(dest, destBackup)
			}

			fmt.Printf("Downloading ian's new version to %s\n", dest)
			if err := ioutil.WriteFile(dest, data, 0755); err != nil {
				os.Rename(destBackup, dest)
				fmt.Printf("%v Failed to update ian\n", color.RedString("Error:"))
				return
			}

			// Removing backup
			os.Remove(destBackup)

			fmt.Printf("ian updated with success to version %s\n", remoteVersion)
		} else {
			fmt.Printf("ian is up to date\n")
		}
	},
}
