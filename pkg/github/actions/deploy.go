/*
Copyright © 2020 Stefan Kürzeder <stefan.kuerzeder@whiteduck.de>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package actions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/whiteducksoftware/azure-arm-action/pkg/azure"
	"github.com/whiteducksoftware/azure-arm-action/pkg/github"
	"github.com/whiteducksoftware/azure-arm-action/pkg/util"
)

func Deploy(ctx context.Context, inputs github.Inputs) error {
	sp, err := azure.GetServicePrincipal(inputs)
	if err != nil {
		return err
	}

	// Load template and parameters if set
	template, err := util.ReadJSON(inputs.TemplateLocation)
	if err != nil {
		return err
	}

	var parameters *map[string]interface{}
	if inputs.ParametersLocation != "" {
		parameters, err = util.ReadJSON(inputs.ParametersLocation)
		if err != nil {
			return err
		}
	}

	// Load authorizer from the service principal
	authorizer, err := azure.GetArmAuthorizerFromServicePrincipal(sp)
	if err != nil {
		return err
	}

	// Load the arm deployments client
	deploymentsClient := azure.GetDeploymentsClient(inputs.SubscriptionID, authorizer)
	inputs.DeploymentName = fmt.Sprintf("%s-%s", inputs.DeploymentName, uuid.New().String())

	// Validate deployment
	logrus.Infof("Validating deployment %s, mode: %s", inputs.DeploymentName, inputs.DeploymentMode)
	validationResult, err := azure.ValidateDeployment(ctx, deploymentsClient, inputs.ResourceGroupName, inputs.DeploymentName, inputs.DeploymentMode, template, parameters)
	if err != nil {
		return err
	}

	if validationResult.StatusCode != http.StatusOK {
		return fmt.Errorf("Template validation failed, %s", validationResult.Status)
	}
	logrus.Info("Validation finished.")

	// Create and wait for completion of the deployment
	logrus.Infof("Creating deployment %s", inputs.DeploymentName)
	resultDeployment, err := azure.CreateDeployment(ctx, deploymentsClient, inputs.ResourceGroupName, inputs.DeploymentName, inputs.DeploymentMode, template, parameters)
	if err != nil {
		return err
	}
	if resultDeployment.StatusCode != http.StatusOK {
		return fmt.Errorf("Template deployment failed, %s", resultDeployment.Status)
	}
	logrus.Info("Template deployment finished.")

	return nil
}
