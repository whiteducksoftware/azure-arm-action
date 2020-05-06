/*
Copyright (c) 2020 white duck Gesellschaft für Softwareentwicklung mbH

This code is licensed under MIT license (see LICENSE for details)
*/
package actions

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/whiteducksoftware/azure-arm-action/pkg/azure"
	"github.com/whiteducksoftware/azure-arm-action/pkg/github"
)

// Authenticate creates and azure authorizer
func Authenticate(inputs github.Inputs) (*autorest.Authorizer, error) {
	var authorizer *autorest.Authorizer

	// Load authorizer from the service principal
	authorizer, err := azure.GetArmAuthorizerFromServicePrincipal(inputs.Credentials)
	if err != nil {
		return authorizer, err
	}

	return authorizer, nil
}
