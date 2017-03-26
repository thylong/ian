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

package projects

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/fatih/color"
)

var httpGet = http.Get

// ErrStatsUnavailable occurs when stats cannot be retrive from github API
var ErrStatsUnavailable = fmt.Errorf("%v %s", color.RedString("Error:"), errors.New("Cannot get stats"))

// ErrJSONPayloadInvalidFormat is returned when the JSON payload format is invalid
var ErrJSONPayloadInvalidFormat = fmt.Errorf("%v %s", color.RedString("Error:"), errors.New("Invalid JSON format"))

// Status makes a GET HTTP query and returns OK if response status is 200
// otherwise ERROR.
func Status(project string, baseURL string, healthEndpoint string) string {
	url := baseURL + healthEndpoint
	resp, err := httpGet(url)
	if err != nil || resp.StatusCode > 300 {
		return fmt.Sprintf("%s: %s", project, color.RedString("ERROR"))
	}

	return fmt.Sprintf("%s: %s", project, color.GreenString("OK"))
}

// Stats makes a GET HTTP query on Github API and returns the stasts.
func Stats(project string, repositoryURL string) (stats map[string]interface{}, err error) {
	resp, err := httpGet(repositoryURL)
	if err != nil || resp.StatusCode > 300 {
		return stats, ErrStatsUnavailable
	}
	content, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err = json.Unmarshal(content, &stats); err != nil {
		return stats, ErrJSONPayloadInvalidFormat
	}

	for k := range stats {
		if strings.HasSuffix(k, "_url") || k == "owner" {
			delete(stats, k)
		}
	}

	return stats, nil
}
