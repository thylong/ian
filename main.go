package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thylong/ian/cmd"
	"github.com/thylong/ian/pkg/log"
)

var version = "undefined"

func init() {
	cmd.RootCmd.AddCommand(versionCmd)
	viper.Set("VERSION", version)
}

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Errorln(err)
		os.Exit(-1)
	}
}

// versionCmd execution displays ian version.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Print the version information.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infoln("ian version:", version)
	},
}
