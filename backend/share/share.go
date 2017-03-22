package share

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/thylong/ian/backend/config"
)

var httpPost = http.Post
var httpGet = http.Get

// Upload to transfer.sh
func Upload(filename string, targetURL string, key string) (string, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		return "", fmt.Errorf("%v Failed to write to buffer", color.RedString("Error:"))
	}

	var fh io.Reader
	if key == "" {
		fh, err = os.Open(filename)
		if err != nil {
			return "", fmt.Errorf("%v Failed to open file", color.RedString("Error:"))
		}
	} else {
		if len(key) < 32 {
			return "", fmt.Errorf("%v The key is too short (less than 32 characters)", color.RedString("Error: "))
		}
		text, _ := ioutil.ReadFile(filename)
		var encrypted []byte
		encrypted, err = EncryptFile(text, []byte(key))
		if err != nil {
			return "", fmt.Errorf("%v Failed to encrypt file", color.RedString("Error:"))
		}
		fh = bytes.NewReader(encrypted)
	}

	// iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return "", err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := httpPost(targetURL, contentType, bodyBuf)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(resp.Status)
	return string(respBody), nil
}

// Download retrieves a distant configuration file
func Download(URL string, configFileName string, key string) error {
	_, err := url.ParseRequestURI(URL)
	if err != nil {
		return fmt.Errorf("%v Sorry, The link you provided is invalid", color.RedString("Error:"))
	}

	resp, err := httpGet(URL)
	if err != nil {
		return fmt.Errorf("%v Sorry, The link you provided is unreachable", color.RedString("Error:"))
	}
	defer resp.Body.Close()

	confFileName := strings.TrimSuffix(configFileName, ".yml")
	if confFilePath, ok := config.ConfigFilesPathes[confFileName]; ok {
		f, err := os.OpenFile(confFilePath, os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		defer f.Close()

		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("%v Cannot read from HTTP response", color.RedString("Error:"))
		}
		if key != "" {
			if content, err = DecryptFile(content, []byte(key)); err != nil {
				return fmt.Errorf("%v Cannot decrypt downloaded file", color.RedString("Error:"))
			}
		}

		if _, err := io.Copy(f, bytes.NewReader(content)); err != nil {
			return fmt.Errorf("%v %s", color.RedString("Error:"), err)
		}
	}
	return nil
}

// GetshareRetrieveFromLinkCmdUsageTemplate returns shareRetrieveFromLinkCmd usage template
func GetshareRetrieveFromLinkCmdUsageTemplate() string {
	return `Usage:{{if .Runnable}}
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
`
}
