package common

import (
	"flag"
	"log"
	"os/exec"
	"strings"
	"math"
	"time"
	"strconv"
)
 const (
	 TotalPortCheckAttempts = 10	 
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
	log.Println("Starting building image...")
	commandPath := Testconfig.DockerfilesHome + "/" + product.Name + "/" + "build.sh"
	logFileName := product.Name + product.Version + RunLogs
	args := " -v " + product.Version + " -r " + product.Provisioning_method + " -q > " + logFileName + " 2>&1"
	command := "bash " + commandPath + args
	_, err := exec.Command("/bin/bash", "-c", command).Output()

	if err == nil {
		log.Println("Successfully built docker image.")
	}
}

func RunImage(product Product) {
	log.Println("Running image...")
	commandPath := Testconfig.DockerfilesHome + "/" + product.Name + "/" + "run.sh"
	logFileName := product.Name + product.Version + RunLogs
	args := " -v " + product.Version + " > " + logFileName + " 2>&1"
	command := "bash " + commandPath + args
	_, err := exec.Command("/bin/bash", "-c", "echo 'n n' | "+command).Output()

	if err == nil {
		log.Println("Successfully ran docker image.")
	}
}

func CheckBuildLogs(productName string, productVersion string) {
	logFileName := productName + productVersion + BuildLogs
	command := "grep -i 'error' " + logFileName
	_, err := exec.Command("/bin/bash", "-c", command).Output()
	if err == nil {
		log.Println("Found errors in docker build logs, please check " + logFileName)
	} else {
		log.Println("No errors were found in docker build logs")
		command := "rm ./" + logFileName
		RunCommandAndGetError(command)
	}
}

func CheckRunLogs(productName string, productVersion string) {
	logFileName := productName + productVersion + RunLogs
	command := "grep -i 'error' " + logFileName
	_, err := exec.Command("/bin/bash", "-c", command).Output()
	if err == nil {
		log.Println("Found errors in docker run logs, please check " + logFileName)
	} else {
		log.Println("No errors were found in docker run logs")
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

func CheckPortWithTimeout(containerIp string, port string) {
	for i := 1; i <= TotalPortCheckAttempts; i++ {
		portCheckCommand := "nc -z -v -w5 " + containerIp + " " + strings.TrimSpace(port)
		err := RunCommandAndGetError(portCheckCommand)

		if err != nil {
			log.Println("Attempt: " + string(i) + ". Port " + port +" is not open in the docker container.")
		}
		
		sleepTime := math.Pow(2, float64(i))
		log.Println("Sleeping for " + strconv.FormatFloat(sleepTime, 'E', -1, 64) + " seconds")		
		time.Sleep(time.Duration(int32(sleepTime)) * time.Second)
	}
}

func CheckWso2CarbonServerStatus(productName string) {
	containerIp := GetDockerContainerIP(productName)
	command := "curl --insecure --write-out %{http_code} --silent --output /dev/null https://" +
		containerIp + ":" + Testconfig.Carbon_Server_Port +
		"/carbon/admin/login.jsp"
	result := RunCommandAndGetOutput(command)
	if result == "200" {
		log.Println("Carbon server is up and running.")
	} else {
		log.Println("Caron server is not running.")
	}
}

func CheckWso2CarbonServerLogs(productName string, productVersion string) {
	log.Println("Checking Carbon server logs for any errors")

	CopyWSO2CarbonLogs(productName, productVersion)
	command := "grep -ir 'error' ./" + productName + productVersion + "logs"
	err := RunCommandAndGetError(command)
	
	if err != nil {
		log.Println("Errors founds in carbon server logs, please check them under " +
			productName + productVersion + "logs")
	} else {
		log.Println("Carbon server logs does not contain any errors")
		command := "rm -rf ./" + productName + productVersion + "logs"
		RunCommandAndGetError(command)
	}
}

func RunCommandAndGetOutput(command string) string {
	out, err := exec.Command("/bin/bash", "-c", command).Output()
	if err != nil {
		log.Fatal("Unable to run command "+command, err)
	}
	return string(out)
}

func RunCommandAndGetError(command string) error {
	_, err := exec.Command("/bin/bash", "-c", command).Output()
	return err
}
