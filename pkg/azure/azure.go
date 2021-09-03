/*
Copyright (c) 2020 white duck Gesellschaft für Softwareentwicklung mbH

This code is licensed under MIT license (see LICENSE for details)
*/
package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/resources/mgmt/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/azure/cli"
)

// SDKAuth represents Azure Sp
type SDKAuth struct {
	ClientID       string `json:"clientId"`
	ClientSecret   string `json:"clientSecret"`
	SubscriptionID string `json:"subscriptionId"`
	TenantID       string `json:"tenantId"`
	ARMEndpointURL string `json:"resourceManagerEndpointUrl"`
}

// GetSdkAuthFromString builds from the cmd flags a ServicePrincipal
func GetSdkAuthFromString(credentials string) (SDKAuth, error) {
	var auth SDKAuth
	err := json.Unmarshal([]byte(credentials), &auth)
	if err != nil {
		return SDKAuth{}, fmt.Errorf("failed to parse the credentials passed, marshal error: %s", err)
	}

	return auth, nil
}

// GetArmAuthorizerFromSdkAuth creates an ARM authorizer from an Sp
func GetArmAuthorizerFromSdkAuth(auth SDKAuth) (*autorest.Authorizer, error) {
	oauthconfig, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, auth.TenantID)
	if err != nil {
		return nil, err
	}

	// If the Resource Manager Endpoint is not set, fallback to the default public cloud endpoint
	if len(auth.ARMEndpointURL) == 0 {
		auth.ARMEndpointURL = azure.PublicCloud.ResourceManagerEndpoint
	}

	token, err := adal.NewServicePrincipalToken(*oauthconfig, auth.ClientID, auth.ClientSecret, auth.ARMEndpointURL)
	if err != nil {
		return nil, err
	}

	// Create authorizer from the bearer token
	var authorizer autorest.Authorizer
	authorizer = autorest.NewBearerAuthorizer(token)

	return &authorizer, nil
}

// GetArmAuthorizerFromSdkAuthJSON creats am ARM authorizer from the passed sdk auth file
func GetArmAuthorizerFromSdkAuthJSON(path string) (*autorest.Authorizer, error) {
	var authorizer autorest.Authorizer

	// Manipulate the AZURE_AUTH_LOCATION var at runtime
	os.Setenv("AZURE_AUTH_LOCATION", path)
	defer os.Unsetenv("AZURE_AUTH_LOCATION")

	authorizer, err := auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
	return &authorizer, err
}

// GetArmAuthorizerFromSdkAuthJSONString creates an ARM authorizer from the sdk auth credentials
func GetArmAuthorizerFromSdkAuthJSONString(credentials string) (*autorest.Authorizer, error) {
	var authorizer autorest.Authorizer

	// create a temporary file, as the sdk credentials need to be read from a file
	tmpFile, err := ioutil.TempFile(os.TempDir(), "azure-sdk-auth-")
	if err != nil {
		return &authorizer, fmt.Errorf("Cannot create temporary sdk auth file: %s", err)
	}
	defer os.Remove(tmpFile.Name())

	text := []byte(credentials)
	if _, err = tmpFile.Write(text); err != nil {
		return &authorizer, fmt.Errorf("Failed to write to temporary sdk auth file: %s", err)
	}
	tmpFile.Close()

	// Manipulate the AZURE_AUTH_LOCATION var at runtime
	os.Setenv("AZURE_AUTH_LOCATION", tmpFile.Name())
	defer os.Unsetenv("AZURE_AUTH_LOCATION")

	authorizer, err = auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)

	return &authorizer, err
}

// GetArmAuthorizerFromEnvironment creates an ARM authorizer from a MSI (AAD Pod Identity)
func GetArmAuthorizerFromEnvironment() (*autorest.Authorizer, error) {
	var authorizer autorest.Authorizer
	authorizer, err := auth.NewAuthorizerFromEnvironment()

	return &authorizer, err
}

// GetArmAuthorizerFromCLI creates an ARM authorizer from the local azure cli
func GetArmAuthorizerFromCLI() (autorest.Authorizer, error) {
	token, err := cli.GetTokenFromCLIWithParams(cli.GetAccessTokenParams{Resource: azure.PublicCloud.ResourceManagerEndpoint})
	if err != nil {
		return nil, err
	}

	adalToken, err := token.ToADALToken()
	if err != nil {
		return nil, err
	}

	return autorest.NewBearerAuthorizer(&adalToken), nil
}

// GetDeploymentsClient takes the azure authorizer and creates an ARM deployments client on the desired subscription
func GetDeploymentsClient(subscriptionID string, authorizer autorest.Authorizer) resources.DeploymentsClient {
	deployClient := resources.NewDeploymentsClient(subscriptionID)
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

// ValidateDeploymentAtSubscriptionScope validates the template deployments and their
// parameters are correct and will produce a successful deployment.GetResource (at subscription scope)
func ValidateDeploymentAtSubscriptionScope(ctx context.Context, deployClient resources.DeploymentsClient, deploymentName string, deploymentMode string, template, params map[string]interface{}) (valid resources.DeploymentValidateResult, err error) {
	return deployClient.ValidateAtSubscriptionScope(ctx,
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
func CreateDeployment(ctx context.Context, deployClient resources.DeploymentsClient, resourceGroupName, deploymentName string, deploymentMode string, template, params map[string]interface{}) (de resources.DeploymentExtended, err error) {
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

// CreateDeploymentAtSubscriptionScope creates a template deployment using the
// referenced JSON files for the template and its parameters (at subscription scope)
func CreateDeploymentAtSubscriptionScope(ctx context.Context, deployClient resources.DeploymentsClient, deploymentName, deploymentMode string, template, params map[string]interface{}) (de resources.DeploymentExtended, err error) {
	future, err := deployClient.CreateOrUpdateAtSubscriptionScope(
		ctx,
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

func GetActiveSubscriptionFromCLI() (string, error) {
	cliCmd := cli.GetAzureCLICommand()
	cliCmd.Args = append(cliCmd.Args, "account", "show", "-o", "json")

	output, err := cliCmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("Invoking Azure CLI failed with the following error: %s", ee.Stderr)
		}

		return "", fmt.Errorf("Invoking Azure CLI failed with the following error: %s", err.Error())
	}

	var data struct {
		SubscriptionID string `json:"id"`
	}
	err = json.Unmarshal(output, &data)
	if err != nil {
		return "", err
	}

	return data.SubscriptionID, err
}
