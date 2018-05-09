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
