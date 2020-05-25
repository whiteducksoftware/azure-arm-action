/*
Copyright (c) 2020 white duck Gesellschaft für Softwareentwicklung mbH

This code is licensed under MIT license (see LICENSE for details)
*/
package actions

import "github.com/mitchellh/mapstructure"

// Output represents a single output of an ARM template
type Output struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// ParseOutputs takes the raw outputs from the azure.DemploymentExtended object
// and converts it to a string Output map
func ParseOutputs(raw interface{}) (map[string]Output, error) {
	if raw == nil {
		return map[string]Output{}, nil
	}

	var outputs map[string]Output
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &outputs,
		TagName: "json",
	})
	if err != nil {
		return map[string]Output{}, err
	}

	err = decoder.Decode(raw)
	if err != nil {
		return map[string]Output{}, err
	}

	return outputs, nil
}