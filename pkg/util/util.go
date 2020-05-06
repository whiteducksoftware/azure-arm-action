/*
Copyright (c) 2020 white duck Gesellschaft für Softwareentwicklung mbH

This code is licensed under MIT license (see LICENSE for details)
*/
package util

import (
	"encoding/json"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

// ReadJSON reads a json file, and unmashals it.
// Very useful for template deployments.
func ReadJSON(path string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Fatalf("failed to read template file: %v\n", err)
	}
	contents := make(map[string]interface{})
	if err := json.Unmarshal(data, &contents); err != nil {
		return nil, err
	}
	return contents, nil
}
