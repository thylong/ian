package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thylong/ian/backend/config"
	"github.com/thylong/ian/backend/log"
	"github.com/thylong/ian/backend/share"
)

var key = ""

var encryptShareCmdParam bool
var decryptShareCmdParam bool

// var shortLinkShareCmdParam bool

func init() {
	RootCmd.AddCommand(shareCmd)

	shareCmd.PersistentFlags().BoolVarP(&encryptShareCmdParam, "encrypt", "e", false, "Encrypt with key before uploading (32 characters minimum)")
	shareCmd.PersistentFlags().BoolVarP(&decryptShareCmdParam, "decrypt", "d", false, "Decrypt config file")
	// shareCmd.PersistentFlags().BoolVarP(&shortLinkShareCmdParam, "bitlink", "b", false, "Get a Bit.ly shorten URL")
	shareRetrieveFromLinkCmd.SetUsageTemplate(share.GetshareRetrieveFromLinkCmdUsageTemplate())

	shareCmd.AddCommand(
		shareConfigCmd,
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
			log.Errorln("It looks like I cannot upload configuration file... :(")
			return
		}
		log.Infoln(link)
	},
}

var shareEnvCmd = &cobra.Command{
	Use:   "env",
	Short: "Share a public link to ian env.yml file",
	Long:  `Share a public link to ian env.yml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if encryptShareCmdParam {
			key = config.GetUserPrivateInput("Enter a secret key (32 characters minimum)")
		}
		link, err := share.Upload(config.ConfigFilesPathes[cmd.Use], "https://transfer.sh/", key)
		if err != nil {
			log.Errorln("It looks like I cannot upload configuration file... :(")
			return
		}
		log.Infoln(link)
	},
}

var shareRetrieveFromLinkCmd = &cobra.Command{
	Use:   "retrieve",
	Short: "Retrieve config from config file link",
	Long:  `Retrieve config from config file link.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Errorln("Not enough argument\n")
			cmd.Usage()
			return
		}

		key := ""
		if decryptShareCmdParam {
			key = config.GetUserInput("Enter the secret key")
		}
		err := share.Download(args[1], args[0], key)
		if err != nil {
			log.Errorf("%s\n", err)
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
