package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/colors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/shipyard-run/shipyard/pkg/config"
	"github.com/shipyard-run/shipyard/pkg/shipyard"
)

var currentClients *shipyard.Clients
var currentConfig *config.Config
var currentEngine *shipyard.Engine

var opt = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestMain(m *testing.M) {
	format := "progress"
	for _, arg := range os.Args[1:] {
		if arg == "-test.v=true" { // go test transforms -v option
			format = "pretty"
			break
		}
	}
	status := godog.RunWithOptions("godog", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format: format,
		Paths:  []string{"features"},
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^the config "([^"]*)"$`, theConfig)
	s.Step(`^I run apply$`, iRunApply)
	s.Step(`^there should be (\d+) container running called "([^"]*)"$`, thereShouldBeContainerRunningCalled)
	s.Step(`^there should be 1 network called "([^"]*)"$`, thereShouldBe1NetworkCalled)
	s.Step(`^a call to "([^"]*)" should result in status (\d+)$`, aCallToShouldResultInStatus)

	s.BeforeScenario(func(interface{}) {
	})

	s.AfterScenario(func(interface{}, error) {
		err := currentEngine.Destroy()
		if err != nil {
			panic(err)
		}
	})
}

func theConfig(arg1 string) error {
	// we are creating this manually rather than using NewWithFolder
	// so we can obtain references to the clients for checks

	var err error
	currentConfig, err = config.New()
	if err != nil {
		return err
	}

	err = config.ParseFolder(arg1, currentConfig)
	if err != nil {
		return err
	}

	err = config.ParseReferences(currentConfig)
	if err != nil {
		return err
	}

	l := hclog.New(&hclog.LoggerOptions{Level: hclog.Debug})

	// create providers
	cc, err := shipyard.GenerateClients(l)
	if err != nil {
		return err
	}

	currentClients = cc
	currentEngine = shipyard.New(currentConfig, cc, l)

	return nil
}

func iRunApply() error {
	return currentEngine.Apply()
}

func thereShouldBeContainerRunningCalled(arg1 int, arg2 string) error {
	// a container can start immediately and then it can crash, this can cause a false positive for the test
	// wait a few seconds to ensure the state does not change
	time.Sleep(5 * time.Second)

	// we need to check this a number of times to make sure it is not just a slow starting container
	for i := 0; i < 10; i++ {
		args := filters.NewArgs()
		args.Add("name", arg2)
		opts := types.ContainerListOptions{Filters: args, All: true}

		cl, err := currentClients.Docker.ContainerList(context.Background(), opts)
		if err != nil {
			return err
		}

		if len(cl) == arg1 {
			// check to see if the container has failed
			if cl[0].State == "exited" {
				return fmt.Errorf("container exited prematurely")
			}

			if cl[0].State == "running" {
				return nil
			}
		}

		// wait a few seconds before trying again
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("Expected %d containers %s", arg1, arg2)
}

func thereShouldBe1NetworkCalled(arg1 string) error {
	args := filters.NewArgs()
	args.Add("name", arg1)
	n, err := currentClients.Docker.NetworkList(context.Background(), types.NetworkListOptions{Filters: args})

	if err != nil {
		return err
	}

	if len(n) != 1 {
		return fmt.Errorf("Expected 1 network called %s to be created", arg1)
	}

	return nil
}

// test making a HTTP call, for testing Ingress
func aCallToShouldResultInStatus(arg1 string, arg2 int) error {
	// try 10 times
	var err error
	for i := 0; i < 10; i++ {
		var resp *http.Response
		resp, err = http.Get(arg1)

		if err == nil && resp.StatusCode == arg2 {
			return nil
		}

		time.Sleep(2 * time.Second)
	}

	return err
}
