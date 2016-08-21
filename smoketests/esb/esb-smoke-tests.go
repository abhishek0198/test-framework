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
	"os"
	"fmt"
	"flag"
	"log"
	"io"
)

var (
	Ip string
	Port string
	User string
	Pass string
	Logger     *log.Logger
)

func main() {
	initializeLogging()

	flag.StringVar(&Ip, "ip", "", "Container IP Address")
	flag.StringVar(&Port, "port", "", "Carbon server port")
	flag.StringVar(&User, "user", "", "Carbon server username")
	flag.StringVar(&Pass, "pass", "", "Carbon server password")

	flag.Parse()

	if checkValidArgs() {
		Logger.Println("INFO Running ESB Smoke tests. IP: " + Ip + " PORT: " + Port)

		result := true

		result = TestProxyServiceCreationAndRemoval() && result

		if result {
			os.Exit(0)
		}
	} else {
		Logger.Println("ERROR Invalid arguments to run smoke test. Check usage using -help")
	}

	os.Exit(100)
}

func checkValidArgs() bool {
	if Ip == "" || Pass == "" || User == "" || Pass == "" {
		return false
	}
	return true
}

func TestProxyServiceCreationAndRemoval() bool {
	result := true
	result = result && LoginToCarbonServer(Ip, Port, User, Pass)
	result = result && CreateProxyService("test", "http://test.com", Ip, Port)
	result = result && DoesProxyServiceExist("test", Ip, Port) && result
	result = result && DeleteProxyService("test", Ip, Port) && result

	return result
}

func initializeLogging() *os.File {
	f, err := os.OpenFile("esb-smoketest-log.txt", os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("ERROR Error in opening file %v \n", err)
	}

	log.SetOutput(f)
	Logger = log.New(io.MultiWriter(f, os.Stdout), "", log.Lshortfile | log.LstdFlags)
	Logger.Println("INFO Logging initialized")
	return f
}