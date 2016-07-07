package main

import "os/exec"
import "log"

// Function to clean docker image by its tag
func CleanDockerImage(tag string) {
    command := "docker history -q " + tag
    _, err := exec.Command("/bin/bash", "-c", command).Output()

    if err == nil {
	log.Println("Removing docker image " + tag)
	_, err := exec.Command("/bin/bash", "-c", "docker rmi " + tag).Output()

        if(err == nil) {
            log.Println("Successfully removed docker image")
        }
    }  
}

// Function to Stop and Remove a docker container 
func StopAndRemoveDockerContainer(productName string) {
    command := "docker ps -a | grep " + productName + " | awk '{print $1}'"
    out, err := exec.Command("/bin/bash", "-c", command).Output()
    
    if(err != nil) {
        log.Println("Error in getting docker container id")
        log.Printf("%s\n", err)
    } else if(string(out) != ""){
	log.Printf("Stopping and removing docker container with id: %s\n", out)

        _, err1 := exec.Command("/bin/bash", "-c", "docker stop " + string(out)).Output()
    	_, err2 := exec.Command("/bin/bash", "-c", "docker rm " + string(out)).Output()

        if(err1 == nil && err2 == nil) {
            log.Println("Successfully stopped and removed docker container")
    	} 
    }
}
