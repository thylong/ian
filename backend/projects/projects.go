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
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/fatih/color"
)

// Status makes a GET HTTP query and returns OK if response status is 200
// otherwise ERROR.
func Status(project string, baseURL string, healthEndpoint string) {
	url := baseURL + healthEndpoint
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v %s.", color.RedString("Error:"), err)
	}
	defer resp.Body.Close()

	if statusCode := resp.StatusCode; statusCode == 200 {
		fmt.Printf("%s : ", project)
		color.Green("OK")
	} else {
		fmt.Printf("%s : ", project)
		color.Red("ERROR")
	}
}

// Stats makes a GET HTTP query on Github API and returns the stasts.
func Stats(project string, repositoryURL string) {
	resp, err := http.Get(repositoryURL)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	content, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var jsonContent map[string]interface{}
	if err = json.Unmarshal(content, &jsonContent); err != nil {
		fmt.Fprintf(os.Stderr, "%v %s.", color.RedString("Error:"), err)
		return
	}

	fmt.Printf("\n%s:", project)
	fmt.Printf("\n- Forks: %v", jsonContent["forks_count"])
	fmt.Printf("\n- Stars: %v", jsonContent["stargazers_count"])
	fmt.Printf("\n- Open Issues: %v", jsonContent["open_issues_count"])
	fmt.Printf("\n- Last update: %v", jsonContent["updated_at"])
}