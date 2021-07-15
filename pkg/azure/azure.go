/*
Copyright (c) 2020 white duck Gesellschaft für Softwareentwicklung mbH

This code is licensed under MIT license (see LICENSE for details)
*/
package azure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/resources/mgmt/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
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
func GetArmAuthorizerFromCLI() (*autorest.Authorizer, error) {
	var authorizer autorest.Authorizer
	authorizer, err := auth.NewAuthorizerFromCLIWithResource(azure.PublicCloud.ResourceManagerEndpoint)

	return &authorizer, err
}

// GetDeploymentsClient takes the azure authorizer and creates an ARM deployments client on the desired subscription
func GetDeploymentsClient(subscriptionID string, authorizer *autorest.Authorizer) resources.DeploymentsClient {
	deployClient := resources.NewDeploymentsClient(subscriptionID)
	deployClient.Authorizer = *authorizer
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
	// This is the path that a developer can set to tell this class what the install path for Azure CLI is.
	const azureCLIPath = "AzureCLIPath"

	// The default install paths are used to find Azure CLI. This is for security, so that any path in the calling program's Path environment is not used to execute Azure CLI.
	azureCLIDefaultPathWindows := fmt.Sprintf("%s\\Microsoft SDKs\\Azure\\CLI2\\wbin; %s\\Microsoft SDKs\\Azure\\CLI2\\wbin", os.Getenv("ProgramFiles(x86)"), os.Getenv("ProgramFiles"))

	// Default path for non-Windows.
	const azureCLIDefaultPath = "/bin:/sbin:/usr/bin:/usr/local/bin"

	// Execute Azure CLI to get subscription id
	var cliCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cliCmd = exec.Command(fmt.Sprintf("%s\\system32\\cmd.exe", os.Getenv("windir")))
		cliCmd.Env = os.Environ()
		cliCmd.Env = append(cliCmd.Env, fmt.Sprintf("PATH=%s;%s", os.Getenv(azureCLIPath), azureCLIDefaultPathWindows))
		cliCmd.Args = append(cliCmd.Args, "/c", "az")
	} else {
		cliCmd = exec.Command("az")
		cliCmd.Env = os.Environ()
		cliCmd.Env = append(cliCmd.Env, fmt.Sprintf("PATH=%s:%s", os.Getenv(azureCLIPath), azureCLIDefaultPath))
	}
	cliCmd.Args = append(cliCmd.Args, "account", "show", "-o", "json")

	var stderr bytes.Buffer
	cliCmd.Stderr = &stderr

	output, err := cliCmd.Output()
	if err != nil {
		return "", fmt.Errorf("Invoking Azure CLI failed with the following error: %s", stderr.String())
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
