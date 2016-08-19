/*
 ------------------------------------------------------------------------

 Copyright 2016 WSO2, Inc. (http://wso2.com)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License

 ------------------------------------------------------------------------
*/

package main

import (
	"encoding/json"
	"github.com/abhishek0198/wso2dockerfiles-test-framework/common"
	"io/ioutil"
	"log"
)

type jsontype struct {
	Testconfig common.TestConfig
}

func parseTestConfig() {
	file, err := ioutil.ReadFile(common.TestConfigPath)
	if err != nil {
		panic("Unable to read test config file." + err.Error())
	}

	log.Println("Parsing test configuration")

	var configObject jsontype
	err = json.Unmarshal(file, &configObject)

	if err != nil {
		panic("Could not parse test config json, please check if its correct.")
	}

	common.Testconfig = configObject.Testconfig
}
