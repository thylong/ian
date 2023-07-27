package cmd

import (
	"os"
	"time"

	"github.com/thylong/ian/pkg/log"

	"github.com/spf13/cobra"
	pm "github.com/thylong/ian/pkg/package-managers"
)

// OSPackageManager is the main package manager used by the current OS.
var OSPackageManager pm.PackageManager

func init() {
	var err error
	OSPackageManager, err = pm.GetOSPackageManager()
	if err != nil {
		log.Errorf("%s\n", err)
		os.Exit(1)
	}
	if !OSPackageManager.IsInstalled() {
		OSPackageManager.Setup()
	}
}

// RootCmd is executed by default (top level).
var RootCmd = &cobra.Command{
	Use:   "ian",
	Short: "Ian is a simple tool to manage your development environment",
	Long:  `Ian is a simple tool to manage your development environment and repositories.`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func spinner() {
	for {
		for _, v := range `-\|/` {
			log.Infof("\rUpdating env... %c", v)
			time.Sleep(100 * time.Millisecond)
		}
	}
}
