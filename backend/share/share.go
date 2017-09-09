package share

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/afero"
	"github.com/thylong/ian/backend/config"
)

// AppFs is a wrapper to OS package
var AppFs = afero.NewOsFs()

var httpPost = http.Post
var httpGet = http.Get

// ErrFailedToOpenFile occurs when trying to open a non-existing config file
var ErrFailedToOpenFile = errors.New("Failed to open file")

// ErrConfiFileMissing occurs when trying to open a non existing config file
var ErrConfiFileMissing = errors.New("Config file doesn't exist")

// ErrLinkUnreachable occurs when a link is unreachable or doe not exist
var ErrLinkUnreachable = errors.New("Sorry, The link you provided is unreachable")

// ErrLinkWrongFormat occurs when a link has a wrong format
var ErrLinkWrongFormat = errors.New("Sorry, The link you provided is invalid")

// Upload to transfer.sh
func Upload(filename string, URL string, key string) (respBody string, err error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, _ := bodyWriter.CreateFormFile("uploadfile", filename)

	var fh io.Reader
	if key == "" {
		fh, err = AppFs.Open(filename)
		if err != nil {
			return "", ErrFailedToOpenFile
		}
	} else {
		text, _ := ioutil.ReadFile(filename)
		var encrypted []byte
		encrypted, err = EncryptFile(text, []byte(key))
		if err != nil {
			return "", err
		}
		fh = bytes.NewReader(encrypted)
	}
	io.Copy(fileWriter, fh)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := httpPost(URL, contentType, bodyBuf)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	return string(respBodyBytes), err
}

// Download retrieves a distant configuration file
func Download(filename string, URL string, key string) (err error) {
	if _, err = url.ParseRequestURI(URL); err != nil {
		return ErrLinkWrongFormat
	}

	resp, err := httpGet(URL)
	if err != nil {
		return ErrLinkUnreachable
	}
	defer resp.Body.Close()

	confFileName := strings.TrimSuffix(filename, ".yml")
	if confFilePath, ok := config.ConfigFilesPathes[confFileName]; ok {
		f, err := AppFs.OpenFile(confFilePath, os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			return ErrConfiFileMissing
		}
		defer f.Close()

		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.New("Cannot read from HTTP response")
		}
		if key != "" {
			if content, err = DecryptFile(content, []byte(key)); err != nil {
				return errors.New("Cannot decrypt downloaded file")
			}
		}
		io.Copy(f, bytes.NewReader(content))
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
