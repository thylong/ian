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
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(newsCmd)
	newsCmd.AddCommand(lastNewsCmd)
	// newsCmd.AddCommand(cityNewsCmd)
}

// newsCmd represents the news command
var newsCmd = &cobra.Command{
	Use:   "news",
	Short: "Retrieve last news from Google News",
	Long:  `We should have to open a browser to know what's going on in the world...`,
}

var lastNewsCmd = &cobra.Command{
	Use:   "last",
	Short: "Retrieve last news",
	Long:  `Retrieve last news.`,
	Run: func(cmd *cobra.Command, args []string) {
		url := "https://news.google.com/news?cf=all&pz=1&ned=us&output=rss"
		displayNews(url)
	},
}

var cityNewsCmd = &cobra.Command{
	Use:   "city",
	Short: "Retrieve news from a city",
	Long:  `Retrieve news from a city.`,
	Run: func(cmd *cobra.Command, args []string) {
		url := "https://news.google.com/news?cf=all&hl=en&pz=1&ned=fr&output=rss&oq=caen"
		displayNews(url)
	},
}

func displayNews(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	type Item struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
	}
	type Items struct {
		XMLName  xml.Name `xml:"channel"`
		ItemList []Item   `xml:"item"`
	}
	type RSS struct {
		XMLName xml.Name `xml:"rss"`
		Items   Items    `xml:"channel"`
	}

	var rss RSS
	err = xml.Unmarshal(content, &rss)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
	for c, item := range rss.Items.ItemList {
		re := regexp.MustCompile("url=(.*)")
		fmt.Println("News :")
		fmt.Printf("\t%d: %s\n", c+1, item.Title)
		fmt.Printf("\t%s\n\n", re.FindString(item.Link))
	}
}
