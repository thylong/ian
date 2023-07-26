// Copyright 2023 Théotime Levêque
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Infof print an info with a colored prefix.
// Errorln formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
func Infof(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(fmt.Sprintf("%s", format), a...)
}

// Infoln formats using a colored prefix prior to its operands and writes to
// standard output. Spaces are always added between operands and a newline is
// appended.
func Infoln(a ...interface{}) {
	for _, b := range a {
		fmt.Fprint(os.Stdout, b)
	}
	fmt.Fprint(os.Stdout, "\n")
}

// Warningf print an info with a colored prefix.
// Errorln formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
func Warningf(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(
		fmt.Sprintf("%v %s", color.YellowString("Warning:"), format),
		a...,
	)
}

// Warningln formats using a colored prefix prior to its operands and writes to
// standard output. Spaces are always added between operands and a newline is
// appended.
func Warningln(a ...interface{}) {
	fmt.Fprint(os.Stderr, color.YellowString("Warning:"), " ")
	for _, b := range a {
		fmt.Fprint(os.Stderr, b)
	}
	fmt.Fprint(os.Stderr, "\n")
}

// Errorf print an error with a colored prefix.
// Errorln formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
func Errorf(format string, a ...interface{}) (n int, err error) {
	fmt.Println(fmt.Sprintf("%v %s", color.RedString("Error:"), format))
	return fmt.Fprintf(
		os.Stderr,
		fmt.Sprintf("%v %s", color.RedString("Error:"), format),
		a...,
	)
}

// Errorln formats using a colored prefix prior to its operands and writes to
// standard error. Spaces are always added between operands and a newline is
// appended.
func Errorln(a ...interface{}) {
	fmt.Fprint(os.Stderr, color.RedString("Error:"), " ")
	for _, b := range a {
		fmt.Fprint(os.Stderr, b)
	}
	fmt.Fprint(os.Stderr, "\n")
}
