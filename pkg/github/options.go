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
package github

import (
	"fmt"
	"reflect"
	"time"

	"github.com/caarlos0/env"
	"github.com/whiteducksoftware/azure-arm-action/pkg/azure"
	"github.com/whiteducksoftware/azure-arm-action/pkg/util"
)

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
	Credentials       azure.ServicePrincipal `env:"INPUT_CREDS"`
	SubscriptionID    string                 `env:"INPUT_SUBSCRIPTIONID"`
	Template          map[string]interface{} `env:"INPUT_TEMPLATELOCATION"`
	Parameters        map[string]interface{} `env:"INPUT_PARAMERTERSLOCATION"`
	ResourceGroupName string                 `env:"INPUT_RESOURCEGROUPNAME"`
	DeploymentName    string                 `env:"INPUT_DEPLOYMENTNAME"`
	DeploymentMode    string                 `env:"INPUT_DEPLOYMENTMODE"`
	Timeout           time.Duration          `env:"INPUT_TIMEOUT" envDefault:"20m"`
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
	reflect.TypeOf(azure.ServicePrincipal{}): wrapGetServicePrincipal,
	reflect.TypeOf(map[string]interface{}{}): wrapReadJSON,
}

func wrapGetServicePrincipal(v string) (interface{}, error) {
	return azure.GetServicePrincipal(v)
}

func wrapReadJSON(v string) (interface{}, error) {
	return util.ReadJSON(v)
}
