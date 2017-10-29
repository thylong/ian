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

package env

import "errors"

// ErrJSONPayloadInvalidFormat is returned when the JSON payload format is invalid
var ErrJSONPayloadInvalidFormat = errors.New("Invalid JSON format")

// ErrOperationNotPermitted is returned when trying create or write without permissions
var ErrOperationNotPermitted = errors.New("Operation not permitted")

// ErrCannotMoveDotfile is returned when trying create or write without permissions
var ErrCannotMoveDotfile = errors.New("Couldn't move dotfile")

// ErrCannotSymlink is returned when trying to create a Symlink and fails
var ErrCannotSymlink = errors.New("Couldn't create symlink")

// ErrCannotInteractWithGit is returned when trying to interact with Git
var ErrCannotInteractWithGit = errors.New("Cannot interact with Git")

// ErrHTTPError is returned when failing to reach an endpoint with HTTP
var ErrHTTPError = errors.New("Cannot reach endpoint")

// ErrDotfilesRepository is returned when failing to stat a repository
var ErrDotfilesRepository = errors.New("dotfiles repository doesn't exists or is not reachable")
