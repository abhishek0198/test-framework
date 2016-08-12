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
	"crypto/tls"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const (
	LoginURL       = "/carbon/admin/login_action.jsp"
	CreateProxyURL = "/carbon/proxyservices/template_pass-through.jsp"
	DeleteProxyURL = "/carbon/service-mgt/delete_service_groups.jsp"
)

var (
	HttpClient            http.Client
	InitializedHttpClient bool
)

/* Function to create http client to be used to send GET/POST requests to wso2 caron server
   The client is only created once and is used for the rest of the commands. It also persists
   cookies so that URLs which require login can be accessed.
*/
func getHttpClient() http.Client {
	if !InitializedHttpClient {
		Logger.Println("Initializing http client")

		options := cookiejar.Options{
			PublicSuffixList: publicsuffix.List,
		}

		jar, err := cookiejar.New(&options)
		check(err)

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		HttpClient = http.Client{Transport: tr, Jar: jar}
		InitializedHttpClient = true
	}
	return HttpClient
}

// Function to login to wso2 carbon server using default credentials
func LoginToCarbonServer(ip string) {
	client := getHttpClient()
	carbonLoginUrl := "https://" + ip + ":" + Testconfig.Carbon_Server_Port + LoginURL
	resp, err := client.PostForm(carbonLoginUrl, url.Values{
		"password": {Testconfig.Carbon_Server_Username},
		"username": {Testconfig.Carbon_Server_Password},
	})
	check(err)
	defer resp.Body.Close()

	htmlcontent, err := ioutil.ReadAll(resp.Body)
	check(err)

	data := string(htmlcontent)
	if resp.StatusCode == 200 && strings.Contains(data, "Signed-in as:") {
		Logger.Println("WSO2 admin user Logged in")
	} else {
		Logger.Println("WSO2 admin user could not login")
	}
}

// Function to test if a proxy service exists in wso2 esb
func DoesProxyServiceExist(serviceName string, ip string) {
	client := getHttpClient()
	url := "https://" + ip + ":" + Testconfig.Carbon_Server_Port + "/services/" + serviceName + "?tryit"
	resp, err := client.Get(url)
	defer resp.Body.Close()
	check(err)

	if resp.StatusCode == 200 {
		Logger.Printf("Proxy service "+serviceName+" exists. Response code: %d\n", resp.StatusCode)
	} else {
		Logger.Printf("Proxy service "+serviceName+" does not exists. Response code: %d\n", resp.StatusCode)
	}
	resp.Body.Close()
}

// Function to create a new proxy service in wso2 esb
func CreateProxyService(proxyName string, targetURL string, ip string) {
	client := getHttpClient()
	args := "?formSubmitted=true&proxyName=" + proxyName + "&targetEndpointMode=url&targetURL=" + targetURL + "&publishWsdlCombo=None&availableTransportsList=https,http,local&trp__https=https&trp__http=http"
	createServiceUrl := "https://" + ip + ":" + Testconfig.Carbon_Server_Port + CreateProxyURL + args
	resp, err := client.Get(createServiceUrl)
	defer resp.Body.Close()

	check(err)
	Logger.Printf("Proxy service creation for service "+proxyName+" completed. Response code: %d\n", resp.StatusCode)
}

// Function to delete an existing proxy service in wso2 esb
func DeleteProxyService(proxyName string, ip string) {
	client := getHttpClient()
	args := "?pageNumber=0&serviceGroups=" + proxyName
	deleteServiceUrl := "https://" + ip + ":" + Testconfig.Carbon_Server_Port + DeleteProxyURL + args
	resp, err := client.Get(deleteServiceUrl)
	defer resp.Body.Close()

	check(err)
	Logger.Printf("Proxy service deletion for service "+proxyName+" completed. Response code: %d\n", resp.StatusCode)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
