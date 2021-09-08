/*
Copyright (c) 2020 white duck Gesellschaft für Softwareentwicklung mbH

This code is licensed under MIT license (see LICENSE for details)
*/
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/whiteducksoftware/azure-arm-action/pkg/github"
	"github.com/whiteducksoftware/azure-arm-action/pkg/github/actions"
	"github.com/whiteducksoftware/golang-utilities/github/actions/io"
)

func init() {
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to info
	if !ok {
		lvl = "info"
	}
	// parse string, this is built-in feature of logrus
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.InfoLevel
	}
	// set global log level
	logrus.SetLevel(ll)
}

func main() {
	opts, err := github.LoadOptions()
	if err != nil {
		logrus.Errorf("failed to load options: %s", err.Error())
		io.WriteError(io.Message{Message: fmt.Sprintf("failed to load options: %s", err.Error())})
		os.Exit(1)
	}

	// read inptus
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	setupInterruptHandler(cancel)

	// Output some information
	if opts.RunningAsAction {
		logrus.Infof("==== Running workflow %s for %s@%s ====", opts.Workflow, opts.Ref, opts.Commit)
	}

	// authenticate
	authorizer, err := actions.Authenticate(opts)
	if err != nil {
		logrus.Errorf("Failed to authenticate with azure: %s", err.Error())
		io.WriteError(io.Message{Message: fmt.Sprintf("Failed to authenticate with azure: %s", err.Error())})
		os.Exit(1)
	}

	// deploy the template
	resultDeployment, err := actions.Deploy(ctx, opts, authorizer)
	if err != nil {
		logrus.Errorf("Failed to deploy the template: %s", err.Error())
		io.WriteError(io.Message{Message: fmt.Sprintf("Failed to deploy the template: %s", err.Error())})
		os.Exit(1)
	}

	// parse the template outputs
	outputs, err := actions.ParseOutputs(resultDeployment.Properties.Outputs)
	if err != nil {
		logrus.Errorf("Failed to parse the template outputs: %s", err.Error())
		io.WriteError(io.Message{Message: fmt.Sprintf("Failed to parse the template outputs: %s", err.Error())})
		os.Exit(1)
	}

	// write the outputs and the deploymentName to our outputs
	io.SetOutput("deploymentName", *resultDeployment.Name)
	for name, output := range outputs {
		io.SetOutput(name, output.Value)
	}

	if opts.RunningAsAction {
		logrus.Info("==== Successfully finished running the workflow ====")
	}
}

func setupInterruptHandler(cancel func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c // wait for the signal
		logrus.Info("Received interrupt, exiting now...")
		cancel()
		os.Exit(1)
	}()
}
