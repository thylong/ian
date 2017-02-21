package cmd

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"os/user"

	yaml "gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd is executed by default (top level).
var RootCmd = &cobra.Command{
	Use:   "ian",
	Short: "Ian is a very simple automation tool for developer with Mac environment",
	Long:  `Ian is a very simple automation tool for developer with Mac environment.`,
}

// ConfigPath represents the path to config directory.
var ConfigPath string

// IanConfigPath represents the path to ian config directory.
var IanConfigPath string

// ConfigFullPath represents the path to config file.
var ConfigFullPath string

// ConfigYaml is used to marshal/unmarshal config file.
type ConfigYaml struct {
	Managers     map[string]map[string]string `json:"managers"`
	Repositories map[string]string            `json:"repositories"`
	Projects     map[string]map[string]string `json:"projects"`
	Setup        map[string][]string          `json:"setup"`
	Packages     map[string]map[string]string `json:"packages"`
}

// Config contains the config content.
var Config ConfigYaml

func init() {
	// Init and keep track of config.
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	ConfigPath = usr.HomeDir + "/.config"
	IanConfigPath = ConfigPath + "/ian/"
	ConfigFullPath = IanConfigPath + "config.yml"
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(IanConfigPath)

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println("Problem with config file: ", err.Error())
	} else {
		configContent, _ := ioutil.ReadFile(ConfigFullPath)
		err = yaml.Unmarshal(configContent, &Config)
		if err != nil {
			log.Warning("Unable to parse config file.")
			return
		}
		viper.WatchConfig()
	}

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

// Execute a command and print output from stdout & stderr.
func executeCommand(termCmd *exec.Cmd) {
	subcmdStds, err := termCmd.CombinedOutput()

	if err != nil {
		log.Info(err)
	}
	fmt.Printf("%s", subcmdStds)
}

// InjectPackagesInTemplate print packages list with usage.
func InjectPackagesInTemplate() string {
	packagesUsages := `Package Commands:
`
	for packageName, packageMeta := range Config.Packages {
		packagesUsages += `  ` + packageName + ` ` + packageMeta["description"] + ` type:` + packageMeta["type"] + "\n"
	}
	return packagesUsages
}
