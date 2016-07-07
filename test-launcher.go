package main

import "log"
import "os"
import "fmt"

const DOCKERFILE_HOME = "/home/abhishek/dev/dockerfiles"

func main() {
    product := CreateProductConfig()
    //productPath := DOCKERFILE_HOME + "/" + product.name

    f, err := os.OpenFile(product.outputFile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        fmt.Printf("Error in opening file %v", err)
    }
    defer f.Close()

    log.SetOutput(f)
    log.Println("Logging initialized")

    log.Println("Running tests for " + product.name + ", " + product.version + " using profile " + product.provisioningMethod)
    //CleanDockerImage(product.name + ":" + product.version)   
    StopAndRemoveDockerContainer(product.name)
}
