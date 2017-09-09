package config

import (
	"os"

	"github.com/thylong/ian/backend/log"
)

// GetPreset returns the content of the preset env.yml
func GetPreset(presetName string) []byte {
	return []byte{}
}

// CreateEnvFileWithPreset write an env file with the selected preset.
func CreateEnvFileWithPreset(preset string) {
	var Envcontent string
	switch preset {
	default:
		log.Errorln("Cannot find preset.")
		return
	case "1":
		Envcontent = GetSoftwareEngineerPreset()
	case "2":
		Envcontent = GetBackendDeveloperPreset()
	case "3":
		Envcontent = GetFrontendDeveloperPreset()
	case "4":
		Envcontent = GetOpsPreset()
	}

	confPath := ConfigFilesPathes["env"]
	f, err := os.OpenFile(confPath, os.O_CREATE|os.O_WRONLY, 0655)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
	defer f.Close()

	if _, err = f.WriteString(Envcontent); err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
}

// GetSoftwareEngineerPreset returns the Software Engineer preset content.
func GetSoftwareEngineerPreset() string {
	return string([]byte(`brew:
- cask
- cmake
- composer
- httpie
- mongodb
- node
- nmap
- php71
- python
- python3
- ruby
- tree
- watch
- wget
cask:
- appcleaner
- atom
- caffeine
- charles
- dash
- docker
- filezilla
- google-chrome
- iterm2
- slack
- spectacle
- virtualbox
- vlc
apt:
- python
- python3
- ruby
- wget
- tree
- watch
- node
- nmap
- mongodb
- slack
- filezilla
- atom
- vlc
- virtualbox
- google-chrome
yum:
- python
- python3
- ruby
- wget
- tree
- watch
- node
- nmap
- mongodb
- slack
- filezilla
- atom
- vlc
- virtualbox
- google-chrome
`))
}

// GetBackendDeveloperPreset returns the Backend Engineer preset content.
func GetBackendDeveloperPreset() string {
	return string([]byte(`brew:
- cask
- cmake
- composer
- httpie
- mongodb
- node
- nmap
- php71
- python
- python3
- ruby
- tree
- watch
- wget
cask:
- appcleaner
- atom
- caffeine
- charles
- dash
- docker
- filezilla
- google-chrome
- iterm2
- slack
- spectacle
- virtualbox
- vlc
apt:
- python
- python3
- ruby
- wget
- tree
- watch
- node
- nmap
- mongodb
- slack
- filezilla
- atom
- vlc
- virtualbox
- google-chrome
yum:
- python
- python3
- ruby
- wget
- tree
- watch
- node
- nmap
- mongodb
- slack
- filezilla
- atom
- vlc
- virtualbox
- google-chrome
`))
}

// GetFrontendDeveloperPreset returns the Frontend Engineer preset content.
func GetFrontendDeveloperPreset() string {
	return string([]byte(`brew:
- cask
- cmake
- composer
- httpie
- mongodb
- node
- nmap
- php71
- python
- python3
- ruby
- tree
- watch
- wget
cask:
- appcleaner
- atom
- caffeine
- charles
- dash
- docker
- filezilla
- google-chrome
- iterm2
- slack
- spectacle
- virtualbox
- vlc
yum:
- python
- python3
- ruby
- wget
- tree
- watch
- node
- nmap
- mongodb
- slack
- filezilla
- atom
- vlc
- virtualbox
- google-chrome
`))
}

// GetOpsPreset returns the Ops Engineer preset content.
func GetOpsPreset() string {
	return string([]byte(`brew:
- cask
- cmake
- composer
- httpie
- mongodb
- node
- nmap
- php71
- python
- python3
- ruby
- tree
- watch
- wget
cask:
- appcleaner
- atom
- caffeine
- charles
- dash
- docker
- filezilla
- google-chrome
- iterm2
- slack
- spectacle
- virtualbox
- vlc
apt:
- python
- python3
- ruby
- wget
- tree
- watch
- node
- nmap
- mongodb
- slack
- filezilla
- atom
- vlc
- virtualbox
- google-chrome
yum:
- python
- python3
- ruby
- wget
- tree
- watch
- node
- nmap
- mongodb
- slack
- filezilla
- atom
- vlc
- virtualbox
- google-chrome
`))
}
