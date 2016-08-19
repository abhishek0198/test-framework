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
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	TotalCheckAttempts = 10
)

func BuildImage(productName string, productVersion string, pMethod string) bool {
	Logger.Println("Starting building docker image for '" + productName + ":" + productVersion + "' ...")
	commandPath := DockerfilesHome + "/" + productName + "/" + "build.sh"
	logFileName := productName + productVersion + RunLogs
	args := " -v " + productVersion + " -r " + pMethod + " -q > " + logFileName + " 2>&1"
	command := "bash " + commandPath + args
	_, err := exec.Command("/bin/bash", "-c", command).Output()

	if err == nil {
		Logger.Println("INFO Successfully built docker image.")
		return true
	} else {
		Logger.Println(GetRedColorFormattedOutputString("ERROR Docker build failure. Error: " + err.Error()))
		Logger.Printf("Build command used: %s", command)
		buf, _ := ioutil.ReadFile(logFileName)
		Logger.Println("BuildLog: " + string(buf))
		os.Remove(logFileName)
		return false
	}
}

func RunImage(productName string, productVersion string) bool {
	Logger.Println("INFO Running docker image for '" + productName + ":" + productVersion + "' ...")
	commandPath := DockerfilesHome + "/" + productName + "/" + "run.sh"
	logFileName := productName + productVersion + RunLogs
	args := " -v " + productVersion + " > " + logFileName + " 2>&1"
	command := "bash " + commandPath + args
	_, err := exec.Command("/bin/bash", "-c", "echo 'n n' | "+command).Output()

	if err == nil {
		Logger.Println("INFO Successfully started the container for " + productName)
		return true
	} else {
		Logger.Println(GetRedColorFormattedOutputString("ERROR Docker run failed. Error: " +  err.Error()))
		Logger.Printf("Run command: %s", command)
		buf, _ := ioutil.ReadFile(logFileName)
		Logger.Println("RunLog: " + string(buf))
		os.Remove(logFileName)
		return false
	}
}

func CheckBuildLogs(productName string, productVersion string) {
	logFileName := productName + productVersion + BuildLogs
	command := "grep -i 'error' " + logFileName
	_, err := exec.Command("/bin/bash", "-c", command).Output()
	if err == nil {
		Logger.Println(GetRedColorFormattedOutputString("WARN Found errors in docker build logs, please check " + logFileName))
	} else {
		Logger.Println("INFO No errors were found in docker build logs")
		command := "rm ./" + logFileName
		RunCommandAndGetError(command)
	}
}

func CheckRunLogs(productName string, productVersion string) {
	logFileName := productName + productVersion + RunLogs
	command := "grep -i 'error' " + logFileName
	_, err := exec.Command("/bin/bash", "-c", command).Output()
	if err == nil {
		Logger.Println(GetRedColorFormattedOutputString("WARN Found errors in docker run logs, please check " + logFileName))
	} else {
		Logger.Println("INFO No errors were found in docker run logs")
		command := "rm ./" + logFileName
		RunCommandAndGetError(command)
	}
}

func CheckExposedPorts(productName string) bool {
	containerIp := GetDockerContainerIP(productName)
	productPath := DockerfilesHome + "/" + productName
	portLine := RunCommandAndGetOutput("grep EXPOSE " + productPath + "/Dockerfile")
	var ports []string
	ports = strings.Split(portLine, " ")

	result := true
	for i := 1; i < len(ports); i++ {
		port := strings.Replace(ports[i], "\n", "", 2)
		result = CheckPortWithTimeout(containerIp, port, i == 1) && result
	}
	return result
}

func CheckPortWithTimeout(containerIp string, port string, applyBackOff bool) bool {
	attempts := 3
	if(applyBackOff) {
		attempts = TotalCheckAttempts
	}
	
	for i := 1; i <= attempts; i++ {
		portCheckCommand := "nc -z -v -w5 " + containerIp + " " + strings.TrimSpace(port)
		err := RunCommandAndGetError(portCheckCommand)

		if err != nil {
			Logger.Println("INFO Attempt: " + strconv.Itoa(i) + ". Port " + port + " is not open in the docker container.")

			sleepTime := 2 * i
			if (i + 1) <= attempts {
				Logger.Println("INFO Sleeping for " + strconv.Itoa(sleepTime) + " seconds")
				time.Sleep(time.Duration(int32(sleepTime)) * time.Second)
			}
		} else {
			Logger.Println("INFO Port " + port + " is open in the docker container.")
			return true
		}
	}

	return false
}

func CheckWso2CarbonServerStatus(productName string) bool {
	client, err := GetHttpClient()
	if err != nil {
		return false
	}

	for i := 1; i <= TotalCheckAttempts; i++ {
		url := "https://" + Testconfig.Docker_Container_Ip + ":" + Testconfig.Carbon_Server_Port + "/carbon/admin/login.jsp"
		resp, err := client.Get(url)

		if err == nil && resp != nil && resp.StatusCode == 200 {
			Logger.Println("INFO Carbon server is up and running.")
			resp.Body.Close()
			return true
		} else {
			errMsg := ""
			if err != nil {
				errMsg = "Error: " +err.Error()
			}
			Logger.Println(GetRedColorFormattedOutputString("WARN Attempt: " + strconv.Itoa(i) + ". Carbon server is not running. " + errMsg))
			sleepTime := 2 * i
			if (i + 1) <= TotalCheckAttempts {
				Logger.Println("INFO Sleeping for " + strconv.Itoa(sleepTime) + " seconds")
				time.Sleep(time.Duration(int32(sleepTime)) * time.Second)
			}
		}
	}
	return false
}

func CheckWso2CarbonServerLogs(productName string, productVersion string) bool {
	Logger.Println("INFO Checking Carbon server logs for any errors")

	CopyWSO2CarbonLogs(productName, productVersion)
	command := "grep -r ' ERROR ' ./" + productName + productVersion + "logs"
	err := RunCommandAndGetError(command)

	if err == nil {
		Logger.Println("WARN Errors founds in carbon server logs, please check them under " +
			productName + productVersion + "logs")
		return false
	} else {
		Logger.Println("INFO Carbon server logs does not contain any errors")
		command := "rm -rf ./" + productName + productVersion + "logs"
		RunCommandAndGetError(command)
		return true
	}
}

func RunCommandAndGetOutput(command string) string {
	out, err := exec.Command("/bin/bash", "-c", command).Output()
	if err != nil {
		Logger.Fatal("ERROR Unable to run command "+command, err)
	}
	return string(out)
}

func RunCommandAndGetError(command string) error {
	_, err := exec.Command("/bin/bash", "-c", command).Output()
	return err
}

func GetRedColorFormattedOutputString(msg string) string {
	return "\x1b[31;1m" + msg + "\x1b[0m"
}