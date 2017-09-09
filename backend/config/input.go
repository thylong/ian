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
