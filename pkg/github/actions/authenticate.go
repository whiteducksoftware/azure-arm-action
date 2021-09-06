/*
Copyright (c) 2020 white duck Gesellschaft für Softwareentwicklung mbH

This code is licensed under MIT license (see LICENSE for details)
*/
package actions

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/whiteducksoftware/azure-arm-action/pkg/github"
)

// Authenticate creates and azure authorizer
func Authenticate(inputs github.Inputs) (autorest.Authorizer, error) {
	// Load authorizer from the service principal
	authorizer, err := inputs.Credentials.GetResourceManagerAuthorizer()
	if err != nil {
		return nil, err
	}

	return authorizer, nil
}
