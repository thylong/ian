package cmd

import (
	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/thylong/ian/pkg/config"
	"github.com/thylong/ian/pkg/env"
	"github.com/thylong/ian/pkg/log"
	pm "github.com/thylong/ian/pkg/package-managers"
)

var allCmdParam bool

func init() {
	RootCmd.AddCommand(envCmd)

	RootCmd.AddCommand(
		envAddCmd,
		envRemoveCmd,
		envShowCmd,
		envSaveCmd,
		envFreezeCmd,
	)
}

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage development environment",
	Long:  `Show details, update and save your development environment.`,
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

var envShowCmd = &cobra.Command{
	Use:   "show",
	Short: "List all packages persisted in Ian configuration",
	Long:  `List all packages persisted in Ian configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		settings := config.Vipers["env"].AllSettings()
		prettySettings, _ := json.MarshalIndent(settings, "", "  ")
		log.Infof("Configuration:\n%s\n}", prettySettings)
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

var envFreezeCmd = &cobra.Command{
	Use:   "freeze",
	Short: "Freeze all packages installed through package managers in Ian env file",
	Long:  `Freeze all packages installed through package managers in Ian env file, they can then be saved to your dotfiles repository or shared.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}
