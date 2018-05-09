package config

import (
	"bufio"
	"bytes"
	"os"
	"strings"

	"github.com/thylong/ian/backend/log"

	"github.com/howeyc/gopass"
)

// GetUserInput ask question and return user input.
func GetUserInput(question string) string {
	log.Infof("%s: ", question)
	reader := bufio.NewReader(os.Stdin)
	if input, _ := reader.ReadString('\n'); input != "\n" && input != "" {
		return string(bytes.TrimSuffix([]byte(input), []byte("\n")))
	}
	return ""
}

// GetUserPrivateInput ask question and return user input (silent stdin).
func GetUserPrivateInput(question string) string {
	log.Infof("%s: ", question)
	pass, _ := gopass.GetPasswd()
	return string(pass)
}

// GetBoolUserInput ask question and return true if the user agreed otherwise false.
func GetBoolUserInput(question string) bool {
	in := GetUserInput(question)

	if strings.ToLower(in) == "y" || strings.ToLower(in) == "yes" || strings.ToLower(in) == "" {
		return true
	}
	return false
}
