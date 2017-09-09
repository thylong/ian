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
	"os"

	"github.com/thylong/ian/backend/log"

	"github.com/spf13/cobra"
	pm "github.com/thylong/ian/backend/package-managers"
)

// OSPackageManager is the main package manager used by the current OS.
var OSPackageManager pm.PackageManager

func init() {
	var err error
	OSPackageManager, err = pm.GetOSPackageManager()
	if err != nil {
		log.Errorf("%s\n", err)
		os.Exit(1)
	}
	if !OSPackageManager.IsInstalled() {
		OSPackageManager.Setup()
	}
}

// RootCmd is executed by default (top level).
var RootCmd = &cobra.Command{
	Use:   "ian",
	Short: "Ian is a simple tool to manage your development environment",
	Long: `Ian is a simple tool to manage your development environment, repositories,
and projects.`,
}
