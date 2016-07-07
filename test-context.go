package main

import "flag"
import "fmt"
import "log"
import "os"

type Product struct {
    name string
    version string 
    provisioningMethod string
    organization string
    platform string
    outputFile string
}

func CreateProductConfig() *Product {
    productName := flag.String("n", "", "product name to test")
    productVersion := flag.String("v", "", "product version to test")
    provisioningMethod := flag.String("r", "default", "provisioning method")
    organization := flag.String("o", "", "Organization name")
    platform := flag.String("p", "default", "Platform to test under")
    
    outputFile := flag.String("f", "./output.txt", "path of the test result file")

    flag.Parse()

    if flag.NFlag() < 2 {
        flag.Usage()
        return nil
    }

    return &Product{*productName, *productVersion, *provisioningMethod, *organization, *platform, *outputFile}
}

func initializeLogging(fileName string) {
    f, err := os.OpenFile(fileName, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        fmt.Printf("Error in opening file %v", err)        
    }
    defer f.Close()

    log.SetOutput(f)
    log.Println("Logging initialized")
}
