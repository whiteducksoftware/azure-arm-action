/*
Copyright (c) 2020 white duck Gesellschaft für Softwareentwicklung mbH

This code is licensed under MIT license (see LICENSE for details)
*/
package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/whiteducksoftware/azure-arm-action/pkg/github"
	"github.com/whiteducksoftware/azure-arm-action/pkg/github/actions"
)

func init() {
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to debug
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
		logrus.Errorf("failed to load options: %s", err)
		os.Exit(1)
	}

	// read inptus
	inputs := opts.Inputs
	ctx, cancel := context.WithTimeout(context.Background(), inputs.Timeout)
	setupInterruptHandler(cancel)

	// Output some information
	githubOptions := opts.GitHub
	if githubOptions.RunningAsAction {
		logrus.Infof("==== Running workflow %s for %s@%s ====", githubOptions.Workflow, githubOptions.Ref, githubOptions.Commit)
	}

	// authenticate
	authorizer, err := actions.Authenticate(inputs)
	if err != nil {
		logrus.Errorf("Failed to authenticate with azure: %s", err)
		os.Exit(1)
	}

	// deploy the template
	resultDeployment, err := actions.Deploy(ctx, inputs, authorizer)
	if err != nil {
		logrus.Errorf("Failed to deploy the template: %s", err)
		os.Exit(1)
	}

	// output the deploymentname
	github.SetOutput("deploymentName", *resultDeployment.Name)
	if githubOptions.RunningAsAction {
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
