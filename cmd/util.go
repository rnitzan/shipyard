package cmd

import (
	"fmt"
	"os"
	"strings"
	"runtime"
	"github.com/hashicorp/go-hclog"
)

// HomeFolder returns the users homefolder this will be $HOME on windows and mac and
// USERPROFILE on windows
func HomeFolder() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}

	return os.Getenv("HOME")
}

// ShipyardHome returns the location of the shipyard
// folder, usually $HOME/.shipyard
func ShipyardHome() string {
	return fmt.Sprintf("%s/.shipyard", HomeFolder())
}

// StateDir returns the location of the shipyard
// state, usually $HOME/.shipyard/state
func StateDir() string {
	return fmt.Sprintf("%s/state", ShipyardHome())
}

// IsLocalFolder tests if the given path is a localfolder and can
// exist in the current filesystem
// TODO make more robust with error messages
// to improve UX
func IsLocalFolder(path string) bool {
	if strings.HasPrefix(path, "/") || strings.HasPrefix(path, "./") {
		// test to see if the folder or file exists
		f, err := os.Open(path)
		if err != nil || f == nil {
			return false
		}

		return true
	}

	return false
}

// GetBlueprintFolder parses a blueprint uri and returns the top level
// blueprint folder
// if the URI is not a blueprint will return an error
func GetBlueprintFolder(blueprint string) (string, error) {
	// get the folder for the blueprint
	parts := strings.Split(blueprint, "//")

	if parts == nil || len(parts) != 2 {
		fmt.Println(parts)
		return "", ErrorInvalidBlueprintURI
	}

	return parts[1], nil
}

func createLogger() hclog.Logger {
	return hclog.New(&hclog.LoggerOptions{Level: hclog.Debug, Color: hclog.AutoColor})
}
