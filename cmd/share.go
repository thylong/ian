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
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/thylong/ian/backend/config"
	"github.com/thylong/ian/backend/share"
)

var key = ""

var encryptShareCmdParam bool
var decryptShareCmdParam bool

// var shortLinkShareCmdParam bool

func init() {
	RootCmd.AddCommand(shareCmd)

	shareCmd.PersistentFlags().BoolVarP(&encryptShareCmdParam, "encrypt", "e", false, "Encrypt with private key before uploading")
	shareCmd.PersistentFlags().BoolVarP(&decryptShareCmdParam, "decrypt", "d", false, "Decrypt config file")
	// shareCmd.PersistentFlags().BoolVarP(&shortLinkShareCmdParam, "bitlink", "b", false, "Get a Bit.ly shorten URL")
	shareRetrieveFromLinkCmd.SetUsageTemplate(share.GetshareRetrieveFromLinkCmdUsageTemplate())

	shareCmd.AddCommand(
		shareConfigCmd,
		shareProjectsCmd,
		shareEnvCmd,
		shareRetrieveFromLinkCmd,
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
			key = config.GetUserInput("Enter a secret key (32 characters minimum)")
		}
		link, err := share.Upload(config.ConfigFilesPathes[cmd.Use], "https://transfer.sh/", key)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v It looks like I cannot upload configuration file... :(.", color.RedString("Error:"))
			return
		}
		fmt.Println(link)
	},
}

var shareRetrieveFromLinkCmd = &cobra.Command{
	Use:   "retrieve",
	Short: "Retrieve config from config file link",
	Long:  `Retrieve config from config file link.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "%v Not enough argument.\n\n", color.RedString("Error:"))
			cmd.Usage()
			os.Exit(1)
		}

		key := ""
		if decryptShareCmdParam {
			key = config.GetUserInput("Enter the secret key")
		}
		err := share.Download(args[0], args[1], key)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v %s.", color.RedString("Error:"), err)
			return
		}
	},
}

var shareAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Share a public link to ian a zip containing all files",
	Long:  `Share a public link to ian env.yml file.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}
