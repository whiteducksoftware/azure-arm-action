/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
package deploy

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/whiteducksoftware/azure-arm-action/pkg/azure"
	"github.com/whiteducksoftware/azure-arm-action/pkg/util"
)

// deployCmd represents the deploy command
var Cmd = &cobra.Command{
	Use: "deploy",
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := cmd.Flags()
		sp, err := azure.GetServicePrincipal(flags)
		if err != nil {
			return err
		}

		// Read flags
		subscriptionID, err := flags.GetString(azure.SubscriptionIdFlagName)
		if err != nil {
			return err
		}

		resourceGroupName, err := flags.GetString(azure.ResourceGroupNameFlagName)
		if err != nil {
			return err
		}

		templateLocation, err := flags.GetString(azure.TemplateLocationFlagName)
		if err != nil {
			return err
		}

		template, err := util.ReadJSON(templateLocation)
		if err != nil {
			return err
		}

		deploymentName, err := flags.GetString(azure.DeploymentNameFlagName)
		if err != nil {
			return err
		}
		deploymentName = fmt.Sprintf("%s-%s", deploymentName, uuid.New().String())

		deploymentMode, err := flags.GetString(azure.DeploymentModeFlagName)
		if err != nil {
			return err
		}
		if deploymentMode == "" {
			deploymentMode = "Incremental"
		}

		parametersLocation, err := flags.GetString(azure.ParametersLocationFlagName)
		if err != nil {
			return err
		}

		var parameters *map[string]interface{}
		if parametersLocation != "" {
			parameters, err = util.ReadJSON(parametersLocation)
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
		deploymentsClient := azure.GetDeploymentsClient(subscriptionID, authorizer)

		// Validate deployment
		logrus.Infof("Validating deployment %s, mode: %s", deploymentName, deploymentMode)
		validationResult, err := azure.ValidateDeployment(cmd.Context(), deploymentsClient, resourceGroupName, deploymentName, deploymentMode, template, parameters)
		if err != nil {
			return err
		}

		if validationResult.StatusCode != http.StatusOK {
			return fmt.Errorf("Template validation failed, %s", validationResult.Status)
		}
		logrus.Info("Validation finished.")

		// Create and wait for completion of the deployment
		logrus.Infof("Creating deployment %s", deploymentName)
		resultDeployment, err := azure.CreateDeployment(cmd.Context(), deploymentsClient, resourceGroupName, deploymentName, deploymentMode, template, parameters)
		if err != nil {
			return err
		}
		if resultDeployment.StatusCode != http.StatusOK {
			return fmt.Errorf("Template deployment failed, %s", resultDeployment.Status)
		}
		logrus.Info("Template deployment finished.")

		return nil
	},
}

func init() {
	// Add flags
	Cmd.PersistentFlags().String(azure.CredentialsFlagName, "", "Credentials")
	Cmd.MarkPersistentFlagRequired(azure.CredentialsFlagName)

	Cmd.PersistentFlags().String(azure.SubscriptionIdFlagName, "", "Subscription")
	Cmd.MarkPersistentFlagRequired(azure.SubscriptionIdFlagName)

	Cmd.PersistentFlags().String(azure.ResourceGroupNameFlagName, "", "ResourceGroupName")
	Cmd.MarkPersistentFlagRequired(azure.ResourceGroupNameFlagName)

	Cmd.PersistentFlags().String(azure.TemplateLocationFlagName, "", "TemplateLocation")
	Cmd.MarkPersistentFlagRequired(azure.TemplateLocationFlagName)

	Cmd.PersistentFlags().String(azure.DeploymentNameFlagName, "", "DeploymentName")
	Cmd.MarkPersistentFlagRequired(azure.DeploymentNameFlagName)

	Cmd.PersistentFlags().String(azure.DeploymentModeFlagName, "", "DeploymentMode")
	Cmd.PersistentFlags().String(azure.ParametersLocationFlagName, "", "Parameters")
}
