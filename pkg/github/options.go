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
	"time"

	"github.com/caarlos0/env"
)

type GitHub struct {
	Workflow   string `env:"GITHUB_WORKFLOW"`
	Action     string `env:"GITHUB_ACTION"`
	Actor      string `env:"GITHUB_ACTOR"`
	Repository string `env:"GITHUB_REPOSITORY"`
	Commit     string `env:"GITHUB_SHA"`
	EventName  string `env:"GITHUB_EVENT_NAME"`
	EventPath  string `env:"GITHUB_EVENT_PATH"`
	Ref        string `env:"GITHUB_REF"`
}

type Inputs struct {
	Credentials        string        `env:"INPUT_CREDS"`
	SubscriptionID     string        `env:"INPUT_SUBSCRIPTIONID"`
	TemplateLocation   string        `env:"INPUT_TEMPLATELOCATION"`
	ParametersLocation string        `env:"INPUT_PARAMERTERSLOCATION"`
	ResourceGroupName  string        `env:"INPUT_RESOURCEGROUPNAME"`
	DeploymentName     string        `env:"INPUT_DEPLOYMENTNAME"`
	DeploymentMode     string        `env:"INPUT_DEPLOYMENTMODE"`
	Timeout            time.Duration `env:"INPUT_TIMEOUT" envDefault:"20m"`
}

type Options struct {
	GitHub GitHub
	Inputs Inputs
}

func LoadOptions() (*Options, error) {
	github := GitHub{}
	if err := env.Parse(&github); err != nil {
		return nil, fmt.Errorf("failed to parse github envrionments: %s", err)
	}

	inputs := Inputs{}
	if err := env.Parse(&inputs); err != nil {
		return nil, fmt.Errorf("failed to parse inputs: %s", err)
	}

	return &Options{
		GitHub: github,
		Inputs: inputs,
	}, nil
}
