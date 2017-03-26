package share

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/afero"
	"github.com/thylong/ian/backend/config"
)

func TestUpload(t *testing.T) {
	AppFs = afero.NewMemMapFs()
	defer func() { AppFs = afero.NewOsFs() }()

	cases := []struct {
		Filename         string
		Key              string
		ResponseMock     string
		ConfigFileExists bool
		ExpectedBody     string
		ExpectedErr      error
	}{
		{"config.yml", "", "https://transfer.sh/zUcVH/config.yml", true, "https://transfer.sh/zUcVH/config.yml\n", nil},
		{"config.yml", "", "https://transfer.sh/zUcVH/config.yml", false, "", ErrFailedToOpenFile},
		{"config.yml", "test", "https://transfer.sh/zUcVH/config.yml", true, "", ErrKeyTooShort},
		{"config.yml", "test", "https://transfer.sh/zUcVH/config.yml", false, "", ErrKeyTooShort},
		{"config.yml", "a very very very very secret key", "https://transfer.sh/zUcVH/config.yml", true, "https://transfer.sh/zUcVH/config.yml\n", nil},
	}
	for _, tc := range cases {
		if tc.ConfigFileExists {
			afero.WriteFile(AppFs, tc.Filename, []byte("test"), 0644)
		}
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, tc.ResponseMock)
		}))
		defer ts.Close()

		respBody, err := Upload(tc.Filename, ts.URL, tc.Key)
		if err != tc.ExpectedErr {
			t.Errorf("Upload func returned wrong error: got %#v want %#v",
				err, tc.ExpectedErr)
		}
		if respBody != tc.ExpectedBody {
			t.Errorf("Upload func returned wrong body: got %#v want %#v",
				respBody, tc.ExpectedBody)
		}

		AppFs.Remove(tc.Filename)
	}
}

func TestDownload(t *testing.T) {
	AppFs = afero.NewMemMapFs()
	defer func() { AppFs = afero.NewOsFs() }()

	cases := []struct {
		Filename         string
		Key              string
		ResponseBodyMock string
		ConfigFileExists bool
		WrongURLFormat   bool
		ExpectedErr      error
	}{
		{"config", "", "https://transfer.sh/zUcVH/config.yml", true, false, nil},
		{"config", "", "https://transfer.sh/zUcVH/config.yml", false, false, ErrConfiFileMissing},
		{"config", "", "", true, true, ErrLinkWrongFormat},
	}
	for _, tc := range cases {
		if tc.ConfigFileExists {
			afero.WriteFile(AppFs, config.ConfigFilesPathes[tc.Filename], []byte("test"), 0644)
		}

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, tc.ResponseBodyMock)
		}))
		defer ts.Close()

		if tc.WrongURLFormat {
			ts.URL = "wrongProtocol/transfer.sh/zUcVH/config.yml"
		}

		err := Download(tc.Filename, ts.URL, tc.Key)
		if err != tc.ExpectedErr {
			t.Errorf("Upload func returned wrong error: got %#v want %#v",
				err, tc.ExpectedErr)
		}

		AppFs.Remove(config.ConfigFilesPathes[tc.Filename])
	}
}

func TestGetshareRetrieveFromLinkCmdUsageTemplate(t *testing.T) {
	if GetshareRetrieveFromLinkCmdUsageTemplate() != `Usage:{{if .Runnable}}
  {{if .HasAvailableFlags}}{{appendIfNotPresent .UseLine "<url> <conf_file> [flags]"}}{{else}}{{.UseLine}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
  {{ .CommandPath}} [command]{{end}}{{if gt .Aliases 0}}
Aliases:
  {{.NameAndAliases}}
{{end}}{{if .HasExample}}
Examples:
{{ .Example }}{{end}}{{if .HasAvailableSubCommands}}
Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}
Flags:
{{.LocalFlags.FlagUsages | trimRightSpace}}{{end}}{{if .HasAvailableInheritedFlags}}
Global Flags:
{{.InheritedFlags.FlagUsages | trimRightSpace}}{{end}}{{if .HasHelpSubCommands}}
Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
` {
		t.Errorf("Template changed")
	}
}
