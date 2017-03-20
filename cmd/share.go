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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/thylong/ian/backend/config"
	"github.com/thylong/ian/backend/share"
)

var key = ""

var encryptShareCmdParam bool

// var shortLinkShareCmdParam bool

func init() {
	RootCmd.AddCommand(shareCmd)

	shareCmd.PersistentFlags().BoolVarP(&encryptShareCmdParam, "encrypt", "s", false, "Encrypt with private key before uploading")
	// shareCmd.PersistentFlags().BoolVarP(&shortLinkShareCmdParam, "bitlink", "b", false, "Get a Bit.ly shorten URL")

	shareSetFromLinkCmd.SetUsageTemplate(share.GetshareSetFromLinkCmdUsageTemplate())

	shareCmd.AddCommand(
		shareConfigCmd,
		shareProjectsCmd,
		shareEnvCmd,
		shareSetFromLinkCmd,
		// shareAllCmd,
	)
}

// shareCmd represents the env command
var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "Share ian configuration",
	Long:  `Share a public link to a single (or multiple) ian configuration file.`,
}

var shareConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Share a public link to ian config.yml file",
	Long:  `Share a public link to ian config.yml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if encryptShareCmdParam {
			key = config.GetUserInput("Enter a secret key: ")
		}
		link, err := share.Upload(config.ConfigFilesPathes[cmd.Use], "https://transfer.sh/", key)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v It looks like I cannot upload configuration file... :(.", color.RedString("Error:"))
			return
		}
		fmt.Println(link)
	},
}

var shareProjectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Share a public link to ian projects.yml file",
	Long:  `Share a public link to ian projects.yml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if encryptShareCmdParam {
			key = config.GetUserInput("Enter a secret key: ")
		}
		link, err := share.Upload(config.ConfigFilesPathes[cmd.Use], "https://transfer.sh/", key)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v It looks like I cannot upload configuration file... :(.", color.RedString("Error:"))
			return
		}
		fmt.Println(link)
	},
}

var shareEnvCmd = &cobra.Command{
	Use:   "env",
	Short: "Share a public link to ian env.yml file",
	Long:  `Share a public link to ian env.yml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if encryptShareCmdParam {
			key = config.GetUserInput("Enter a secret key: ")
		}
		link, err := share.Upload(config.ConfigFilesPathes[cmd.Use], "https://transfer.sh/", key)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v It looks like I cannot upload configuration file... :(.", color.RedString("Error:"))
			return
		}
		fmt.Println(link)
	},
}

var shareSetFromLinkCmd = &cobra.Command{
	Use:   "set",
	Short: "Set config from config file link",
	Long:  `Set config from config file link.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "%v Not enough argument.\n\n", color.RedString("Error:"))
			cmd.Usage()
			os.Exit(1)
		}

		_, err := url.ParseRequestURI(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v Sorry, The link you provided is invalid.", color.RedString("Error:"))
		}

		resp, err := http.Get(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v Sorry, The link you provided is unreachable.", color.RedString("Error:"))
		}
		defer resp.Body.Close()

		if strings.HasSuffix(string(strings.TrimSuffix(args[0], "/")), "_e") {
			// DecryptFile
		}

		confFileName := strings.TrimSuffix(args[1], ".yml")
		fmt.Print(confFileName)
		if confFilePath, ok := config.ConfigFilesPathes[confFileName]; ok {
			f, err := os.OpenFile(confFilePath, os.O_TRUNC|os.O_WRONLY, 0600)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer f.Close()

			if _, err := io.Copy(f, resp.Body); err != nil {
				fmt.Fprintf(os.Stderr, "%v %s", color.RedString("Error:"), err)
				os.Exit(1)
			}
			return
		}
		fmt.Fprintf(os.Stderr, "%v Sorry, something went wrong.", color.RedString("Error:"))
	},
}

var shareAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Share a public link to ian a zip containing all files",
	Long:  `Share a public link to ian env.yml file.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}
