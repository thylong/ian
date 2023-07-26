package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/thylong/ian/pkg/log"

	"github.com/mitchellh/ioprogress"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var httpGet = http.Get
var version string
var ianLastCommit = "https://api.github.com/repos/thylong/ian/commits/master"
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
			log.Errorln("Cannot retrieve ian's last release infos.")
			return
		}
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Cannot retrieve ian's last release infos\n")
			return
		}
		defer resp.Body.Close()

		var lastReleaseContent map[string]interface{}
		err = json.Unmarshal(content, &lastReleaseContent)
		if err != nil {
			log.Errorln("Cannot retrieve ian's last version")
			return
		}
		remoteVersion := lastReleaseContent["tag_name"]

		if localVersion != remoteVersion {
			log.Infoln("ian version outdated...")
			downloadURL := lastReleaseContent["assets"].([]interface{})[0].(map[string]interface{})["browser_download_url"]
			resp, err = httpGet(downloadURL.(string))
			if err != nil {
				log.Errorln("Cannot download ian's last version")
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
				log.Errorln("ian's last version seems broken")
				return
			}
			dest, err := os.Executable()
			if err != nil {
				log.Errorln("%v ian's last version seems not executable")
				return
			}

			// Move the old version to a backup path that we can recover from
			// in case the upgrade fails
			destBackup := dest + ".bak"
			if _, err := os.Stat(dest); err == nil {
				os.Rename(dest, destBackup)
			}

			log.Infof("Downloading ian's new version to %s\n", dest)
			if err := ioutil.WriteFile(dest, data, 0755); err != nil {
				os.Rename(destBackup, dest)
				log.Errorln("Failed to update ian")
				return
			}

			// Removing backup
			os.Remove(destBackup)

			log.Infof("ian updated with success to version %s\n", remoteVersion)
		} else {
			log.Infoln("ian is up to date")
		}
	},
}
