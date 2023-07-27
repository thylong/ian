package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thylong/ian/pkg/env"
	"github.com/thylong/ian/pkg/log"
	pm "github.com/thylong/ian/pkg/package-managers"
)

func init() {
	RootCmd.AddCommand(
		envAddCmd,
		envRemoveCmd,
		envSaveCmd,
	)
}

var envAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new package(s) to ian configuration",
	Long:  `Add new package(s) to ian env.yml.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Errorln("Not enough argument")
			cmd.Usage()
			return
		}

		packageManagerName := args[0]
		packages := args[1:]

		if !pm.IsSupportedPackageManager(packageManagerName) {
			log.Errorf("Package Manager %s is not supported\n", packageManagerName)
			return
		}

		env.AddPackagesToEnvFile(packageManagerName, packages)
		log.Infof("Package(s) added to %s list\n", packageManagerName)
	},
}

var envRemoveCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove package(s) to ian configuration",
	Long:  `Remove package(s) to ian env.yml.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Errorln("Not enough argument")
			cmd.Usage()
			return
		}

		packageManagerName := args[0]
		packages := args[1:]

		if !pm.IsSupportedPackageManager(packageManagerName) {
			log.Errorf("Package Manager %s is not supported\n", packageManagerName)
			return
		}

		env.RemovePackagesFromEnvFile(packageManagerName, packages)
		log.Infof("Package(s) removed to %s list\n", args[0])
	},
}

var envSaveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save current configuration files to the dotfiles repository",
	Long:  `Save current configuration files to the dotfiles repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := env.Save([]string{}); err != nil {
			log.Errorf("Save command failed: %s\n", err)
		}
	},
}
