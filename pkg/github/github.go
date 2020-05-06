/*
Copyright (c) 2020 white duck Gesellschaft für Softwareentwicklung mbH

This code is licensed under MIT license (see LICENSE for details)
*/
package github

import "fmt"

// SetOutput can be used to set outputs of your action
func SetOutput(name string, value string) {
	fmt.Printf("::set-output name=%s::%s", name, value)
}
