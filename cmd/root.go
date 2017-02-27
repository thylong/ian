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
	"github.com/spf13/cobra"
	pm "github.com/thylong/ian/backend/package-managers"
)

// OSPackageManager is the main package manager used by the current OS.
var OSPackageManager pm.PackageManager

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

{{end}}{{ if .HasAvailableLocalFlags}}
Flags:
{{.LocalFlags.FlagUsages | trimRightSpace}}{{end}}{{ if .HasAvailableInheritedFlags}}
Global Flags:
{{.InheritedFlags.FlagUsages | trimRightSpace}}{{end}}{{if .HasHelpSubCommands}}
Additional help topics:{{range .Commands}}{{if .IsHelpCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableSubCommands }}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}`)))
}

// RootCmd is executed by default (top level).
var RootCmd = &cobra.Command{
	Use:   "ian",
	Short: "Ian is a very simple automation tool for developer with Mac environment",
	Long:  `Ian is a very simple automation tool for developer with Mac environment.`,
}
