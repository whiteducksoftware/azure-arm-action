/*
Copyright (c) 2020 white duck Gesellschaft für Softwareentwicklung mbH

This code is licensed under MIT license (see LICENSE for details)
*/
package actions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/resources/mgmt/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/whiteducksoftware/azure-arm-action/pkg/azure"
	"github.com/whiteducksoftware/azure-arm-action/pkg/github"
	"github.com/whiteducksoftware/azure-arm-action/pkg/util"
)

// Deploy takes our inputs and initaite and
// waits for completion of the arm template deployment
func Deploy(ctx context.Context, inputs github.Inputs, authorizer autorest.Authorizer) (resources.DeploymentExtended, error) {
	var err error

	// Load the arm deployments client
	deploymentsClient := azure.GetDeploymentsClient(inputs.Credentials.SubscriptionID, authorizer)
	u := uuid.New().String()
	logrus.Infof("Creating deployment %s with uuid %s -> %s-%s, mode: %s", inputs.DeploymentName, u, inputs.DeploymentName, u, inputs.DeploymentMode)
	inputs.DeploymentName = fmt.Sprintf("%s-%s", inputs.DeploymentName, u)

	// Build our final parameters
	parameter := util.MergeParameters(inputs.Parameters, inputs.OverrideParameters)

	// Validate deployment
	logrus.Infof("Validating deployment %s", inputs.DeploymentName)
	var validationResult resources.DeploymentValidateResult

	// check whenether we need to deploy at resource group or subscription scope
	if len(inputs.ResourceGroupName) > 0 {
		validationResult, err = azure.ValidateDeployment(ctx, deploymentsClient, inputs.ResourceGroupName, inputs.DeploymentName, inputs.DeploymentMode, inputs.Template, parameter)
	} else {
		validationResult, err = azure.ValidateDeploymentAtSubscriptionScope(ctx, deploymentsClient, inputs.DeploymentName, inputs.DeploymentMode, inputs.Template, parameter)
	}

	if err != nil {
		return resources.DeploymentExtended{}, err
	}

	if validationResult.StatusCode != http.StatusOK {
		return resources.DeploymentExtended{}, fmt.Errorf("%s, %s", validationResult.Status, *validationResult.Error.Message)
	}
	logrus.Info("Validation finished.")

	// Create and wait for completion of the deployment
	logrus.Infof("Creating deployment %s", inputs.DeploymentName)
	var resultDeployment resources.DeploymentExtended

	// check whenether we need to deploy at resource group or subscription scope
	if len(inputs.ResourceGroupName) > 0 {
		resultDeployment, err = azure.CreateDeployment(ctx, deploymentsClient, inputs.ResourceGroupName, inputs.DeploymentName, inputs.DeploymentMode, inputs.Template, parameter)
	} else {
		resultDeployment, err = azure.CreateDeploymentAtSubscriptionScope(ctx, deploymentsClient, inputs.DeploymentName, inputs.DeploymentMode, inputs.Template, parameter)
	}

	if err != nil {
		return resources.DeploymentExtended{}, err
	}

	// verify the status
	if resultDeployment.StatusCode != http.StatusOK {
		return resources.DeploymentExtended{}, fmt.Errorf("%s", resultDeployment.Status)
	}
	logrus.Info("Template deployment finished.")

	return resultDeployment, nil
}
