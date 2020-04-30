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
package azure

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Flag name constants
const (
	CredentialsFlagName        string = "credentials"
	SubscriptionIdFlagName     string = "subscriptionId"
	ResourceGroupNameFlagName  string = "resourceGroupName"
	TemplateLocationFlagName   string = "templateLocation"
	DeploymentModeFlagName     string = "deploymentMode"
	DeploymentNameFlagName     string = "deploymentName"
	ParametersLocationFlagName string = "parametersLocation"
)

// ServicePrincipal represents Azure Sp
type ServicePrincipal struct {
	Tenant   string
	AppID    string
	Password string
}

// GetServicePrincipalFromFlags builds from the cmd flags a ServicePrincipal
func GetServicePrincipal(credentials string) (*ServicePrincipal, error) {
	var sp ServicePrincipal
	err := json.Unmarshal([]byte(credentials), &sp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the credentials passed, marshal error: %s", err)
	}

	return &sp, nil
}

// GetArmAuthorizerFromServicePrincipal creates an ARM authorizer from an Sp
func GetArmAuthorizerFromServicePrincipal(sp *ServicePrincipal) (*autorest.BearerAuthorizer, error) {
	oauthconfig, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, sp.Tenant)
	if err != nil {
		return nil, err
	}

	token, err := adal.NewServicePrincipalToken(*oauthconfig, sp.AppID, sp.Password, azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		return nil, err
	}

	return autorest.NewBearerAuthorizer(token), nil
}

// GetArmAuthorizerFromMSI creates an ARM authorizer from a MSI (AAD Pod Identity)
func GetArmAuthorizerFromMSI() (autorest.Authorizer, error) {
	return auth.NewAuthorizerFromEnvironment()
}

func GetDeploymentsClient(subscriptionId string, authorizer autorest.Authorizer) resources.DeploymentsClient {
	deployClient := resources.NewDeploymentsClient(subscriptionId)
	deployClient.Authorizer = authorizer
	return deployClient
}

// ValidateDeployment validates the template deployments and their
// parameters are correct and will produce a successful deployment.GetResource
func ValidateDeployment(ctx context.Context, deployClient resources.DeploymentsClient, resourceGroupName, deploymentName string, deploymentMode string, template, params map[string]interface{}) (valid resources.DeploymentValidateResult, err error) {
	return deployClient.Validate(ctx,
		resourceGroupName,
		deploymentName,
		resources.Deployment{
			Properties: &resources.DeploymentProperties{
				Template:   template,
				Parameters: params,
				Mode:       resources.DeploymentMode(deploymentMode),
			},
		})
}

// CreateDeployment creates a template deployment using the
// referenced JSON files for the template and its parameters
func CreateDeployment(ctx context.Context, deployClient resources.DeploymentsClient, resourceGroupName string, deploymentName string, deploymentMode string, template, params map[string]interface{}) (de resources.DeploymentExtended, err error) {
	future, err := deployClient.CreateOrUpdate(
		ctx,
		resourceGroupName,
		deploymentName,
		resources.Deployment{
			Properties: &resources.DeploymentProperties{
				Template:   template,
				Parameters: params,
				Mode:       resources.DeploymentMode(deploymentMode),
			},
		},
	)
	if err != nil {
		return de, fmt.Errorf("cannot create deployment: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, deployClient.Client)
	if err != nil {
		return de, fmt.Errorf("cannot get the create deployment future respone: %v", err)
	}

	return future.Result(deployClient)
}
