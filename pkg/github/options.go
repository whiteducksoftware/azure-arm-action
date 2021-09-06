/*
Copyright (c) 2020 white duck Gesellschaft für Softwareentwicklung mbH

This code is licensed under MIT license (see LICENSE for details)
*/
package github

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"

	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
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
	Credentials        azure.SDKAuth `env:"INPUT_CREDS"`
	Template           template      `env:"INPUT_TEMPLATELOCATION"`
	Parameters         parameters    `env:"INPUT_PARAMETERS"`
	OverrideParameters parameters    `env:"INPUT_OVERRIDEPARAMETERS"`
	ResourceGroupName  string        `env:"INPUT_RESOURCEGROUPNAME"`
	DeploymentName     string        `env:"INPUT_DEPLOYMENTNAME"`
	DeploymentMode     string        `env:"INPUT_DEPLOYMENTMODE"`
	Timeout            time.Duration `env:"INPUT_TIMEOUT" envDefault:"20m"`
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
	reflect.TypeOf(parameters{}):    wrapReadParameters,
}

func wrapGetServicePrincipal(v string) (interface{}, error) {
	return azure.GetSdkAuthFromString(v)
}

func wrapReadJSON(v string) (interface{}, error) {
	logrus.Debugf("Parsing raw json %s", v)
	return util.ReadJSON(v)
}

func wrapReadParameters(v string) (interface{}, error) {
	isJSONInput := strings.HasSuffix(v, ".json") // Todo: This check should be more resilient
	if isJSONInput == true {                     // Check if we are dealing with a path to a json file or raw parameters
		return wrapReadParametersJSON(v)
	}

	return wrapReadRawParameters(v)
}

func wrapReadParametersJSON(v string) (interface{}, error) {
	logrus.Debugf("Parsing parameter json %s", v)
	json, err := util.ReadJSON(v)
	if err != nil {
		return json, err
	}

	// Check if the parameters are wrapped (https://github.com/Azure/azure-sdk-for-go/issues/9283)
	paramters, ok := json["parameters"]
	if ok {
		return paramters, nil
	}

	return json, nil
}

func wrapReadRawParameters(v string) (interface{}, error) {
	parameter := make(map[string]interface{})

	// https://stackoverflow.com/questions/44277222/golang-regular-expression-for-parsing-key-value-pair-into-a-string-map
	lastQuote := rune(0)
	f := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)
		}
	}

	// splitting string by space/newline but considering quoted section
	pairs := strings.FieldsFunc(v, f)

	for _, keyValuePair := range pairs {
		keyValue := strings.SplitN(keyValuePair, "=", 2)
		if len(keyValue) != 2 {
			return nil, fmt.Errorf("Found invalid pair, expected KEY=VALUE got %s", keyValuePair)
		}

		// remove all unicode quotation marks (Todo: should really all be removed?)
		value := strings.Map(func(r rune) rune {
			if unicode.In(r, unicode.Quotation_Mark) {
				return -1
			}
			return r
		}, keyValue[1])

		parameter[keyValue[0]] = make(map[string]string)
		parameter[keyValue[0]].(map[string]string)["value"] = strings.TrimSpace(value)
	}

	return parameter, nil
}
