package main

import (
	"log"
	"encoding/json"
	"io/ioutil"
	"common"
)

type jsontype struct {
    Testconfig common.TestConfig
}


func ParseTestConfig() {
	file, err := ioutil.ReadFile(common.TestConfigPath)
    if (err != nil) {
    	panic("Unable to read test config file." + err.Error())
    }
    
    log.Println("Parsing test configuration")
    
    var configObject jsontype
    json.Unmarshal(file, &configObject)
    common.Testconfig = configObject.Testconfig
}

