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
	"io"
	"log"
	"os"
	"strconv"
	"time"
	"os/exec"
	"strings"
)

const (
	DOCKERFILES_HOME = "DOCKERFILES_HOME"
	PUPPET_HOME = "PUPPET_HOME"
	PUPPET_PROVISIONING = "puppet"
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
	common.Logger.Println("INFO Tests completed in " + totalTime.String())
}

func setDockerFilesHome() {
	if (os.Getenv(DOCKERFILES_HOME) == "") {
		fmt.Println("ERROR DOCKERFILES_HOME is not set. Please set the environment variable before running test")
		os.Exit(1)
	}
	common.DockerfilesHome = os.Getenv(DOCKERFILES_HOME)
}

func initializeLogging() *os.File {
	f, err := os.OpenFile(common.Testconfig.Output_file, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("ERROR Error in opening file %v \n", err)
	}

	log.SetOutput(f)
	common.Logger = log.New(io.MultiWriter(f, os.Stdout), "", log.Lshortfile | log.LstdFlags)
	common.Logger.Println("INFO Logging initialized")
	return f
}

func runTests() {
	var ProductsToTest []common.Product
	testRunToResultMap := make(map[string]string)

	ProductsToTest = common.Testconfig.Products

	for _, product := range ProductsToTest {
		enabled, err := strconv.ParseBool(product.Enabled)

		if err != nil {
			log.Println("ERROR Could not parse the 'enabled' field of test config for test - " + product.Name + ":" + product.Version)
			continue
		}

		if enabled {
			var provisioningMethods []string
			provisioningMethods = product.Provisioning_method

			for _, pMethod := range provisioningMethods {
				doTestSetup(product.Name, product.Version)

				test := product.Name + ":" + product.Version + ":" + pMethod + ":" + product.Platform
				testResult := runSingleTest(product, pMethod)

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
		if result == "Pass" {
			common.Logger.Println("INFO Test:'", test, "', Result:", result)
		} else {
			common.Logger.Println(common.GetRedColorFormattedOutputString("INFO Test:'" + test + "', Result:" + result))
		}
	}
	common.Logger.Println("===========================================================================================")
}

func runSingleTest(product common.Product, pMethod string) bool {
	common.Logger.Println("INFO Running tests for " + product.Name + ", " + product.Version +
		", using " + pMethod + " provisioning" +
		", under " + product.Platform + " platform.")

	if strings.Compare(pMethod, PUPPET_PROVISIONING) == 0 {
		if !isPuppetHomeDefined() {
			return false
		}
	}

	if common.DoesDockerImageExist(product.Name + ":" + product.Version) {
		common.Logger.Println("WARN There is an existing Docker image for this product. Skipping test!")
		return false
	}

	buildResult := common.BuildImage(product.Name, product.Version, pMethod)
	if buildResult {
		common.CheckBuildLogs(product.Name, product.Version)
		if !common.DoesDockerImageExist(product.Name + ":" + product.Version) {
			common.Logger.Println("WARN Docker build was not successful. Skipping test!")
			return false
		}
	} else {
		return false
	}

	runResult := common.RunImage(product.Name, product.Version)
	if runResult {
		common.CheckRunLogs(product.Name, product.Version)

		if !common.IsDockerContainerRunning(product.Name) {
			common.Logger.Println("WARN Docker container is not running. Skipping test")
			return false
		}
	} else {
		return false
	}

	result := true
	result = common.CheckExposedPorts(product.Name) && result
	result = common.CheckWso2CarbonServerStatus() && result
	result = common.CheckWso2CarbonServerLogs(product.Name, product.Version) && result

	// Run smoke tests for this product (if available)
	if product.Smoke_Test_Path!= "" {
		result =  runSmokeTest(product) && result
	}

	common.Logger.Println("INFO Test completed for " + product.Name + ", " + product.Version + ".")
	return result
}

func runSmokeTest(product common.Product) bool {
	common.Logger.Println("INFO Lunching smoke test for " + product.Name + " using " + product.Smoke_Test_Path)
	command := product.Smoke_Test_Path + " -ip " + common.GetDockerContainerIP(product.Name) + " -port " + common.Testconfig.Carbon_Server_Port + " -user " + common.Testconfig.Carbon_Server_Username + " -pass " + common.Testconfig.Carbon_Server_Password

	runResult := exec.Command("/bin/bash", "-c", command)

	_, err := runResult.Output()
	if err != nil {
		common.Logger.Println("ERROR Error in running smoke tests for " + product.Name + ". Error: " + err.Error())
		return false
	}
	out := ""
	r := runResult.ProcessState.Success()
	if r {
		out = "Successful"
	} else {
		out = "Failed"
	}
	common.Logger.Println("INFO Smoke test result for " + product.Name + ": " + out)
	return r
}

func doTestSetup(name string, version string) {
	common.Logger.Println("INFO Starting test setup")
	cleanExistingDockerImage(name, version)
	common.Logger.Println("INFO Completed test setup")
}

func doTestCleanup(name string, version string) {
	common.Logger.Println("INFO Starting test clean up")
	cleanExistingDockerImage(name, version)
	common.Logger.Println("INFO Completed test clean up")
}

func cleanExistingDockerImage(name string, version string) {
	common.StopAndRemoveDockerContainer(name)
	common.CleanDockerImage(name + ":" + version)
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

func isPuppetHomeDefined() bool {
	if (os.Getenv(PUPPET_HOME) == "") {
		fmt.Println("ERROR PUPPET_HOME is not set. Please set the environment variable before running test")
		return false
	}
	return true
}