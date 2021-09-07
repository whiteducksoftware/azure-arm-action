/*
Copyright (c) 2020 white duck Gesellschaft für Softwareentwicklung mbH

This code is licensed under MIT license (see LICENSE for details)
*/
package actions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/whiteducksoftware/azure-arm-action/pkg/github"
	"github.com/whiteducksoftware/azure-arm-action/pkg/util"
	"github.com/whiteducksoftware/golang-utilities/azure/resources/deployments"
)

// Deploy takes our inputs and initaite and
// waits for completion of the arm template deployment
func Deploy(ctx context.Context, options github.Options, authorizer autorest.Authorizer) (resources.DeploymentExtended, error) {
	var err error

	// Load the arm deployments client
	deploymentsClient := deployments.GetClientWithBaseUri(options.Credentials.ARMEndpointURL, options.Credentials.SubscriptionID, authorizer)
	u := uuid.New().String()
	logrus.Infof("Creating deployment %s with uuid %s -> %s-%s, mode: %s", options.DeploymentName, u, options.DeploymentName, u, options.DeploymentMode)
	options.DeploymentName = fmt.Sprintf("%s-%s", options.DeploymentName, u)

	// Build our final parameters
	parameter := util.MergeParameters(options.Parameters, options.OverrideParameters)

	// Validate deployment
	logrus.Infof("Validating deployment %s", options.DeploymentName)

	var validationResult resources.DeploymentValidateResult
	if len(options.ResourceGroupName) > 0 {
		validationResult, err = deployments.Validate(ctx, deploymentsClient, options.ResourceGroupName, options.DeploymentName, options.DeploymentMode, options.Template, parameter)
	} else if len(options.ManagementGroupId) > 0 {
		validationResult, err = deployments.ValidateAtManagementGroupScope(ctx, deploymentsClient, options.ManagementGroupId, options.DeploymentName, options.DeploymentMode, options.Template, parameter)
	} else {
		validationResult, err = deployments.ValidateAtSubscriptionScope(ctx, deploymentsClient, options.DeploymentName, options.DeploymentMode, options.Template, parameter)
	}

	if err != nil {
		return resources.DeploymentExtended{}, err
	}

	if validationResult.StatusCode != http.StatusOK {
		return resources.DeploymentExtended{}, fmt.Errorf("%s, %s", validationResult.Status, *validationResult.Error.Message)
	}
	logrus.Info("Validation finished.")

	// Create and wait for completion of the deployment
	logrus.Infof("Creating deployment %s", options.DeploymentName)

	var resultDeployment resources.DeploymentExtended
	if len(options.ResourceGroupName) > 0 {
		resultDeployment, err = deployments.Create(ctx, deploymentsClient, options.ResourceGroupName, options.DeploymentName, options.DeploymentMode, options.Template, parameter)
	} else if len(options.ManagementGroupId) > 0 {
		resultDeployment, err = deployments.CreateAtManagementGroupScope(ctx, deploymentsClient, options.ManagementGroupId, options.DeploymentName, options.DeploymentMode, options.Template, parameter)
	} else {
		resultDeployment, err = deployments.CreateAtSubscriptionScope(ctx, deploymentsClient, options.DeploymentName, options.DeploymentMode, options.Template, parameter)
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
