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
	"common"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	startTime := time.Now()

	ParseTestConfig()
	f := InitializeLogging(common.Testconfig.Output_file)
	defer f.Close()

	var ProductsToTest []common.Product
	ProductsToTest = common.Testconfig.Wso2_products
	for _, product := range ProductsToTest {
		RunTest(product)
	}

	totalTime := time.Now().Sub(startTime)
	log.Println("Tests completed in " + totalTime.String())
}

func InitializeLogging(outputFile string) *os.File {
	f, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Error in opening file %v", err)
	}

	log.SetOutput(f)
	log.Println("Logging initialized")
	return f
}

func RunTest(product common.Product) {
	enabled, err := strconv.ParseBool(product.Enabled)

	if err != nil {
		panic("Could not parse the 'Enabled' field of test config")
	}

	if enabled {
		log.Println("Running tests for " + product.Name + ", " + product.Version +
			" using profile " + product.Provisioning_method +
			", using platform: " + product.Platform)

		common.CleanDockerImage(product.Name + ":" + product.Version)
		common.StopAndRemoveDockerContainer(product.Name)
		common.BuildImage(product)
		common.CheckBuildLogs(product.Name, product.Version)
		common.RunImage(product)
		common.CheckRunLogs(product.Name, product.Version)
		common.CheckExposedPorts(product.Name)
		common.CheckWso2CarbonServerStatus(product.Name)
		common.CheckWso2CarbonServerLogs(product.Name, product.Version)
		
		// Do cleanup
		common.StopAndRemoveDockerContainer(product.Name)
		common.CleanDockerImage(product.Name + ":" + product.Version)
		
		// Reset globals for next product test run
		common.ResetTestSpecificVariables()
		log.Println("Test completed for " + product.Name + ", " + product.Version + ". \n\n")
	}
}
