package main

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/whiteducksoftware/azure-arm-action/pkg/github"
	"github.com/whiteducksoftware/azure-arm-action/pkg/github/actions"
)

var (
	opts             *github.Options
	authorizer       *autorest.Authorizer
	deploymentResult resources.DeploymentExtended
)

func TestLoadOptions(t *testing.T) {
	var err error
	opts, err = github.LoadOptions()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestAuthentication(t *testing.T) {
	var err error
	authorizer, err = actions.Authenticate(opts.Inputs)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestDeploy(t *testing.T) {
	var err error
	deploymentResult, err = actions.Deploy(context.Background(), opts.Inputs, authorizer)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestParseOutputs(t *testing.T) {
	outputs, err := actions.ParseOutputs(deploymentResult.Properties.Outputs)
	if err != nil {
		t.Error(err.Error())
	}

	if len(outputs) != 2 {
		t.Errorf("Got invalid count of outputs, expected 2 got %d", len(outputs))
	}

	// Test output key location
	value, ok := outputs["location"]
	if !ok {
		t.Errorf("Test key is missing in the outputs, exptected the key location to be present")
	}

	if value.Value != "westeurope" {
		t.Errorf("Got invalid value for location key, expected %s got %s", "westeurope", value.Value)
	}

	// Test output key containername
	value, ok = outputs["containerName"]
	if !ok {
		t.Errorf("Test key is missing in the outputs, exptected the key containerName to be present")
	}

	if value.Value != "github-action" {
		t.Errorf("Got invalid value for location key, expected %s got %s", "github-action", value.Value)
	}
}
