package azure

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/spf13/pflag"
)

// Flag name constants
const (
	CredentialsFlagName       string = "credentials"
	ResourceGroupNameFlagName string = "resourceGroupName"
	TemplateLocationFlagName  string = "templateLocation"
	DeploymentModeFlagName    string = "deploymentMode"
	DeploymentNameFlagName    string = "deploymentName"
	ParametersFlagName        string = "parameters"
)

// ServicePrincipal represents Azure Sp
type ServicePrincipal struct {
	Tenant   string
	AppID    string
	Password string
}

// GetServicePrincipalFromFlags builds from the cmd flags a ServicePrincipal
func GetServicePrincipal(flags *pflag.FlagSet) (*ServicePrincipal, error) {
	credentials, err := flags.GetString(CredentialsFlagName)
	if err != nil {
		return nil, err
	}

	var sp ServicePrincipal
	err = json.Unmarshal([]byte(credentials), &sp)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse the credentials passed:\n\tMarshal Error: %s", err)
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
