package main

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/resources/mgmt/resources"
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

	if len(outputs) != 3 {
		t.Errorf("Got invalid count of outputs, expected  got %d", len(outputs))
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

	// This also tests if the override did work
	if value.Value != "github-action-overriden" {
		t.Errorf("Got invalid value for containerName key, expected %s got %s", "github-action-overriden", value.Value)
	}

	// Test output key containername
	value, ok = outputs["connectionString"]
	if !ok {
		t.Errorf("Test key is missing in the outputs, exptected the key containerName to be present")
	}

	// This also tests if the override did work
	var expectedConnectionString = "Server=tcp:test.database.windows.net;Database=test;User ID=test;Password=test;Trusted_Connection=False;Encrypt=True;"
	if value.Value != expectedConnectionString {
		t.Errorf("Got invalid value for connectionString key, expected %s got %s", expectedConnectionString, value.Value)
	}
}
