package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thylong/ian/backend/config"
	pm "github.com/thylong/ian/backend/packages/managers"
)

// OSPackageManager is the main package manager used by the current OS.
var OSPackageManager pm.PackageManager

// RootCmd is executed by default (top level).
var RootCmd = &cobra.Command{
	Use:   "ian",
	Short: "Ian is a very simple automation tool for developer with Mac environment",
	Long:  `Ian is a very simple automation tool for developer with Mac environment.`,
}

func init() {
	OSPackageManager = pm.GetOSPackageManager()

	RootCmd.SetUsageTemplate(string([]byte(`Usage:{{if .Runnable}}
  {{if .HasAvailableFlags}}{{appendIfNotPresent .UseLine "[flags]"}}{{else}}{{.UseLine}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
  {{ .CommandPath}} [command]{{end}}{{if gt .Aliases 0}}
Aliases:
  {{.NameAndAliases}}
{{end}}{{if .HasExample}}
Examples:
{{ .Example }}{{end}}{{ if .HasAvailableSubCommands}}

Default Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

{{InjectPackagesInTemplate}}
{{end}}{{ if .HasAvailableLocalFlags}}
Flags:
{{.LocalFlags.FlagUsages | trimRightSpace}}{{end}}{{ if .HasAvailableInheritedFlags}}
Global Flags:
{{.InheritedFlags.FlagUsages | trimRightSpace}}{{end}}{{if .HasHelpSubCommands}}
Additional help topics:{{range .Commands}}{{if .IsHelpCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableSubCommands }}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}`)))
	cobra.AddTemplateFunc("InjectPackagesInTemplate", InjectPackagesInTemplate)
}

// InjectPackagesInTemplate prints packages list with usage.
func InjectPackagesInTemplate() string {
	packagesUsages := `Package Commands:
`
	for packageName, packageMeta := range config.Config.Packages {
		packagesUsages += `  ` + packageName + ` ` + packageMeta["description"] + ` type:` + packageMeta["type"] + "\n"
	}
	return packagesUsages
}
