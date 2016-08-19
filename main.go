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
	"fmt"
	"github.com/abhishek0198/wso2dockerfiles-test-framework/common"
	"github.com/abhishek0198/wso2dockerfiles-test-framework/smoketests"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	startTime := time.Now()

	parseTestConfig()
	setDockerFilesHome()
	

	f := initializeLogging()
	defer f.Close()

	if shouldContinue() {
		runTests()
	}

	totalTime := time.Now().Sub(startTime)
	common.Logger.Println("Tests completed in " + totalTime.String())
}
func setDockerFilesHome() {
	if(os.Getenv("DOCKERFILES_HOME") == "") {
		panic("DOCKERFILES_HOME is not set. Please set the environment variable before running test")
	}
	common.DockerfilesHome = os.Getenv("DOCKERFILES_HOME")
}

func initializeLogging() *os.File {
	f, err := os.OpenFile(common.Testconfig.Output_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Error in opening file %v", err)
	}

	log.SetOutput(f)
	common.Logger = log.New(io.MultiWriter(f, os.Stdout), "", log.Lshortfile|log.LstdFlags)
	common.Logger.Println("Logging initialized")
	return f
}

func runTests() {
	var ProductsToTest []common.Product
	testRunToResultMap := make(map[string]string)

	ProductsToTest = common.Testconfig.Wso2_products

	for _, product := range ProductsToTest {
		enabled, err := strconv.ParseBool(product.Enabled)

		if err != nil {
			panic("Could not parse the 'Enabled' field of test config for test - " + product.Name + ":" + product.Version)
			continue
		}

		if enabled {
			var provisioningMethods []string
			provisioningMethods = product.Provisioning_method

			for _, pMethod := range provisioningMethods {
				doTestSetup(product.Name, product.Version)

				test := product.Name + ":" + product.Version + ":" + pMethod + ":" + product.Platform
				testResult := runSingleTest(product.Name, product.Version, pMethod, product.Platform)

				if testResult {
					testRunToResultMap[test] = "Pass"
				} else {
					testRunToResultMap[test] = "Fail"
				}

				doTestCleanup(product.Name, product.Version)
				common.Logger.Println()
			}
		}
	}

	common.Logger.Println("============================  TEST RUN RESULT  ============================================")
	for test, result := range testRunToResultMap {
		common.Logger.Println("Test:'", test, "', Result:", result)
	}
	common.Logger.Println("===========================================================================================")
}

func runSingleTest(name string, version string, pMethod string, platform string) bool {
	common.Logger.Println("Running tests for " + name + ", " + version +
		", using " + pMethod + " provisioning" +
		", under " + platform + " platform.")

	if common.DoesDockerImageExist(name + ":" + version) {
		common.Logger.Println("There is an existing Docker image for this product. Skipping test!")
		return false
	}

	buildResult := common.BuildImage(name, version, pMethod)
	if buildResult {
		common.CheckBuildLogs(name, version)
		if !common.DoesDockerImageExist(name + ":" + version) {
			common.Logger.Println("Docker build was not successful. Skipping test!")
			return false
		}
	} else {
		return false
	}

	runResult := common.RunImage(name, version)
	if runResult {
		common.CheckRunLogs(name, version)

		if !common.IsDockerContainerRunning(name) {
			common.Logger.Println("Docker container is not running. Skipping test")
			return false
		}
	} else {
		return false
	}

	result := true
	result = common.CheckExposedPorts(name) && result
	result = common.CheckWso2CarbonServerStatus(name) && result
	result = common.CheckWso2CarbonServerLogs(name, version) && result

	// Run smoke tests for this product (if available)
	smoketests.RunSmokeTest(name)

	common.Logger.Println("Test completed for " + name + ", " + version + ".")
	return result
}

func doTestSetup(name string, version string) {
	common.Logger.Println("Starting test setup")
	common.StopAndRemoveDockerContainer(name)
	common.CleanDockerImage(name + ":" + version)
	common.Logger.Println("Completed test setup")
}

func doTestCleanup(name string, version string) {
	common.Logger.Println("Starting test clean up")
	common.StopAndRemoveDockerContainer(name)
	common.CleanDockerImage(name + ":" + version)
	common.Logger.Println("Completed test clean up")
}

// Function to check all the preconditions that should meet before we can run the tests. Namely:
// - Docker daemon should be up
// - Add any preconditions here
func shouldContinue() bool {
	if common.IsDockerDaemonRunning() {
		return true
	}
	return false
}
