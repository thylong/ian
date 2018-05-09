package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thylong/ian/backend/env"
	"github.com/thylong/ian/backend/log"
)

func init() {
	RootCmd.AddCommand(restore)
}

// restore represents the setup command
var restore = &cobra.Command{
	Use:   "restore",
	Short: "Restore ian configuration",
	Long:  `Ian requires you to be able to interact with Github through Git CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		env.Restore(OSPackageManager)

		log.Infoln("Great! You're ready to start using Ian.")
	},
}
