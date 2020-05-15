/*
Copyright (c) 2020 white duck Gesellschaft für Softwareentwicklung mbH

This code is licensed under MIT license (see LICENSE for details)
*/
package github

import (
	"fmt"
	"reflect"
	"time"

	"github.com/caarlos0/env"
	"github.com/whiteducksoftware/azure-arm-action/pkg/azure"
	"github.com/whiteducksoftware/azure-arm-action/pkg/util"
)

// Wrapper types to define how we need to parse the input json
type template map[string]interface{}
type parameters map[string]interface{}

// GitHub represents the inputs which github provides us on default
type GitHub struct {
	Workflow        string `env:"GITHUB_WORKFLOW"`
	Action          string `env:"GITHUB_ACTION"`
	Actor           string `env:"GITHUB_ACTOR"`
	Repository      string `env:"GITHUB_REPOSITORY"`
	Commit          string `env:"GITHUB_SHA"`
	EventName       string `env:"GITHUB_EVENT_NAME"`
	EventPath       string `env:"GITHUB_EVENT_PATH"`
	Ref             string `env:"GITHUB_REF"`
	RunningAsAction bool   `env:"GITHUB_ACTIONS" envDefault:"false"`
}

// Inputs represents our custom inputs for the action
type Inputs struct {
	Credentials       azure.SDKAuth `env:"INPUT_CREDS"`
	Template          template      `env:"INPUT_TEMPLATELOCATION"`
	Parameters        parameters    `env:"INPUT_PARAMETERSLOCATION"`
	ResourceGroupName string        `env:"INPUT_RESOURCEGROUPNAME"`
	DeploymentName    string        `env:"INPUT_DEPLOYMENTNAME"`
	DeploymentMode    string        `env:"INPUT_DEPLOYMENTMODE"`
	Timeout           time.Duration `env:"INPUT_TIMEOUT" envDefault:"20m"`
}

// Options is a combined struct of all inputs
type Options struct {
	GitHub GitHub
	Inputs Inputs
}

// LoadOptions parses the environment vars and reads github options and our custom inputs
func LoadOptions() (*Options, error) {
	github := GitHub{}
	if err := env.Parse(&github); err != nil {
		return nil, fmt.Errorf("failed to parse github envrionments: %s", err)
	}

	inputs := Inputs{}
	err := env.ParseWithFuncs(&inputs, customTypeParser)
	if err != nil {
		return nil, fmt.Errorf("failed to parse inputs: %s", err)
	}

	return &Options{
		GitHub: github,
		Inputs: inputs,
	}, nil
}

// custom type parser
var customTypeParser = map[reflect.Type]env.ParserFunc{
	reflect.TypeOf(azure.SDKAuth{}): wrapGetServicePrincipal,
	reflect.TypeOf(template{}):      wrapReadJSON,
	reflect.TypeOf(parameters{}):    wrapReadParamtersJSON,
}

func wrapGetServicePrincipal(v string) (interface{}, error) {
	return azure.GetSdkAuthFromString(v)
}

func wrapReadJSON(v string) (interface{}, error) {
	return util.ReadJSON(v)
}

func wrapReadParamtersJSON(v string) (interface{}, error) {
	json, err := util.ReadJSON(v)
	if err != nil {
		return json, err
	}

	// Check if the parameters are wrapped
	paramters, ok := json["parameters"]
	if ok {
		return paramters, err
	}

	return json, err
}
