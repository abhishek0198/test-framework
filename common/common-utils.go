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

package common

import (
	"flag"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	TotalCheckAttempts = 10
)

func CreateProductConfig() *Product {
	productName := flag.String("n", "", "product name to test")
	productVersion := flag.String("v", "", "product version to test")
	provisioningMethod := flag.String("r", "default", "provisioning method")
	organization := flag.String("o", "", "Organization name")
	platform := flag.String("p", "default", "Platform to test under")

	flag.Parse()

	if flag.NFlag() < 2 {
		flag.Usage()
		return nil
	}

	return &Product{"true", *productName, *productVersion, *provisioningMethod, *organization, *platform}
}

func BuildImage(product Product) {
	Logger.Println("Starting building image...")
	commandPath := Testconfig.DockerfilesHome + "/" + product.Name + "/" + "build.sh"
	logFileName := product.Name + product.Version + RunLogs
	args := " -v " + product.Version + " -r " + product.Provisioning_method + " -q > " + logFileName + " 2>&1"
	command := "bash " + commandPath + args
	_, err := exec.Command("/bin/bash", "-c", command).Output()

	if err == nil {
		Logger.Println("Successfully built docker image.")
	}
}

func RunImage(product Product) {
	Logger.Println("Running image...")
	commandPath := Testconfig.DockerfilesHome + "/" + product.Name + "/" + "run.sh"
	logFileName := product.Name + product.Version + RunLogs
	args := " -v " + product.Version + " > " + logFileName + " 2>&1"
	command := "bash " + commandPath + args
	_, err := exec.Command("/bin/bash", "-c", "echo 'n n' | "+command).Output()

	if err == nil {
		Logger.Println("Successfully ran docker image.")
	}
}

func CheckBuildLogs(productName string, productVersion string) {
	logFileName := productName + productVersion + BuildLogs
	command := "grep -i 'error' " + logFileName
	_, err := exec.Command("/bin/bash", "-c", command).Output()
	if err == nil {
		Logger.Println("Found errors in docker build logs, please check " + logFileName)
	} else {
		Logger.Println("No errors were found in docker build logs")
		command := "rm ./" + logFileName
		RunCommandAndGetError(command)
	}
}

func CheckRunLogs(productName string, productVersion string) {
	logFileName := productName + productVersion + RunLogs
	command := "grep -i 'error' " + logFileName
	_, err := exec.Command("/bin/bash", "-c", command).Output()
	if err == nil {
		Logger.Println("Found errors in docker run logs, please check " + logFileName)
	} else {
		Logger.Println("No errors were found in docker run logs")
		command := "rm ./" + logFileName
		RunCommandAndGetError(command)
	}
}

func CheckExposedPorts(productName string) {
	containerIp := GetDockerContainerIP(productName)
	productPath := Testconfig.DockerfilesHome + "/" + productName
	portLine := RunCommandAndGetOutput("grep EXPOSE " + productPath + "/Dockerfile")
	var ports []string
	ports = strings.Split(portLine, " ")

	for i := 1; i < len(ports); i++ {
		port := strings.Replace(ports[i], "\n", "", 2)
		CheckPortWithTimeout(containerIp, port)
	}
}

func CheckPortWithTimeout(containerIp string, port string) bool {
	for i := 1; i <= TotalCheckAttempts; i++ {
		portCheckCommand := "nc -z -v -w5 " + containerIp + " " + strings.TrimSpace(port)
		err := RunCommandAndGetError(portCheckCommand)

		if err != nil {
			Logger.Println("Attempt: " + strconv.Itoa(i) + ". Port " + port + " is not open in the docker container.")

			sleepTime := 2 * i
			Logger.Println("Sleeping for " + strconv.Itoa(sleepTime) + " seconds")
			time.Sleep(time.Duration(int32(sleepTime)) * time.Second)
		} else {
			Logger.Println("Port " + port + " is open in the docker container.")
			return true
		}
	}

	return false
}

func CheckWso2CarbonServerStatus(productName string) {
	containerIp := GetDockerContainerIP(productName)
	command := "curl --insecure --write-out %{http_code} --silent --output /dev/null https://" +
		containerIp + ":" + Testconfig.Carbon_Server_Port +
		"/carbon/admin/login.jsp"

	for i := 1; i <= TotalCheckAttempts; i++ {
		result := RunCommandAndGetOutput(command)
		if result == "200" {
			Logger.Println("Carbon server is up and running.")
			return
		} else {
			Logger.Println("Attempt: " + strconv.Itoa(i) + "Carbon server is not running.")
			sleepTime := 2 * i
			Logger.Println("Sleeping for " + strconv.Itoa(sleepTime) + " seconds")
			time.Sleep(time.Duration(int32(sleepTime)) * time.Second)
		}
	}
}

func CheckWso2CarbonServerLogs(productName string, productVersion string) {
	Logger.Println("Checking Carbon server logs for any errors")

	CopyWSO2CarbonLogs(productName, productVersion)
	command := "grep -ir 'error' ./" + productName + productVersion + "logs"
	err := RunCommandAndGetError(command)

	if err == nil {
		Logger.Println("Errors founds in carbon server logs, please check them under " +
			productName + productVersion + "logs")
	} else {
		Logger.Println("Carbon server logs does not contain any errors")
		command := "rm -rf ./" + productName + productVersion + "logs"
		RunCommandAndGetError(command)
	}
}

func RunCommandAndGetOutput(command string) string {
	out, err := exec.Command("/bin/bash", "-c", command).Output()
	if err != nil {
		Logger.Fatal("Unable to run command "+command, err)
	}
	return string(out)
}

func RunCommandAndGetError(command string) error {
	_, err := exec.Command("/bin/bash", "-c", command).Output()
	return err
}
