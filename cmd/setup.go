package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thylong/ian/backend/env"
	"github.com/thylong/ian/backend/log"
)

func init() {
	RootCmd.AddCommand(setupCmd)
}

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up ian configuration",
	Long:  `Ian requires you to be able to interact with Github through Git CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		env.Setup(OSPackageManager)

		log.Infoln("Great! You're ready to start using Ian.")
	},
}
