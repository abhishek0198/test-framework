package common

import (
	"log"
	"os/exec"
	"strings"
)

// Function to clean docker image by its tag
func CleanDockerImage(tag string) {
	command := "docker history -q " + tag
	_, err := exec.Command("/bin/bash", "-c", command).Output()

	if err == nil {
		log.Println("Removing docker image " + tag)
		_, err := exec.Command("/bin/bash", "-c", "docker rmi "+tag).Output()

		if err == nil {
			log.Println("Successfully removed docker image")
		}
	}
}

// Function to Stop and Remove a docker container
func StopAndRemoveDockerContainer(productName string) {
	command := "docker ps -a | grep " + productName + " | awk '{print $1}'"
	out, err := exec.Command("/bin/bash", "-c", command).Output()

	if err != nil {
		log.Println("Error in getting docker container id")
		log.Printf("%s\n", err)
	} else if string(out) != "" {
		log.Printf("Stopping and removing docker container with id: %s", out)

		_, err1 := exec.Command("/bin/bash", "-c", "docker stop "+string(out)).Output()
		_, err2 := exec.Command("/bin/bash", "-c", "docker rm "+string(out)).Output()

		if err1 == nil && err2 == nil {
			log.Println("Successfully stopped and removed docker container")
		}
	}
}

// Function to copy WSO2 Carbon serer logs from running container to local
func CopyWSO2CarbonLogs(productName string, productVersion string) {
	containerID := GetDockerContainerID(productName)
	containerIP := GetDockerContainerIP(productName)
	command := "docker cp " + containerID + ":/mnt/" + containerIP + "/" + productName + "-" +
		productVersion + "/repository/logs/ ./" + productName + productVersion + "logs"
	err := RunCommandAndGetError(command)
	if err != nil {
		log.Fatal("Unable to copy carbon server logs from docker container")
	} else {
		log.Println("Successfully copied carbon server logs from docker container")
	}
}

// Function to get docker container id for WSO2 product
func GetDockerContainerID(productName string) string {
	if DockerContainerID == "" {
		command := "docker ps | grep " + productName + " | awk '{print $1}'"
		out, err := exec.Command("/bin/bash", "-c", command).Output()
		if err != nil {
			message := "Error in getting Docker container id. Command: " + command + ". Error: " + err.Error()
			log.Println(message)
			panic(message)
		}

		DockerContainerID = strings.Replace(string(out), "\n", "", 1)
		log.Println("Setting docker ID " + DockerContainerID)
	}

	return DockerContainerID
}

// Function to get docker container IP addrses
func GetDockerContainerIP(productName string) string {
	if DockerContainerIP == "" {
		containerId := GetDockerContainerID(productName)

		command := "docker inspect --format '{{ .NetworkSettings.IPAddress }}' " + containerId
		out, err := exec.Command("/bin/bash", "-c", command).Output()
		if err != nil {
			message := "Error in getting docker container IP. Command: " + command + ", Error: " + err.Error()
			log.Println(message)
			panic(message)
		}
		DockerContainerIP = strings.Replace(string(out), "\n", "", 1)
		log.Println("Setting docker IP " + DockerContainerIP)
	}

	return DockerContainerIP
}
