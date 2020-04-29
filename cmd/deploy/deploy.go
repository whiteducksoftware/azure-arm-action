/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
package deploy

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/whiteducksoftware/azure-arm-action/pkg/azure"
)

// deployCmd represents the deploy command
var Cmd = &cobra.Command{
	Use:   "deploy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := cmd.Flags()
		sp, err := azure.GetServicePrincipal(flags)
		if err != nil {
			return err
		}
		logrus.Info(sp.AppID)

		/*
			authorizer, err := azure.GetArmAuthorizerFromServicePrincipal(sp)
			if err != nil {
				return err
			}
		*/

		return nil
	},
}

func init() {
	// Add flags
	Cmd.PersistentFlags().String(azure.CredentialsFlagName, "", "Credentials")
	Cmd.MarkFlagRequired(azure.CredentialsFlagName)
}
