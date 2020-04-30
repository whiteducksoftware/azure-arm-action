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
package main

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/whiteducksoftware/azure-arm-action/pkg/github"
	"github.com/whiteducksoftware/azure-arm-action/pkg/github/actions"
)

func main() {
	opts, err := github.LoadOptions()
	if err != nil {
		logrus.Errorf("failed to load options: %s", err)
		os.Exit(1)
	}

	// read inptus
	inputs := opts.Inputs
	ctx, cancel := context.WithTimeout(context.Background(), inputs.Timeout)
	defer cancel()

	// deploy the template
	resultDeployment, err := actions.Deploy(ctx, inputs)
	if err != nil {
		logrus.Errorf("failed to load options: %s", err)
		os.Exit(1)
	}

	// output the deploymentname
	github.SetOutput("deploymentName", *resultDeployment.Name)
}
