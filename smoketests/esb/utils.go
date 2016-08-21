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
	"crypto/tls"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

const (
	LoginURL = "/carbon/admin/login_action.jsp"
	CreateProxyURL = "/carbon/proxyservices/template_pass-through.jsp"
	DeleteProxyURL = "/carbon/service-mgt/delete_service_groups.jsp"
)

var (
	HttpClient http.Client
	InitializedHttpClient bool
)

/* Function to create http client to be used to send GET/POST requests to wso2 caron server
   The client is only created once and is used for the rest of the commands. It also persists
   cookies so that URLs which require login can be accessed.
*/
func GetHttpClient() (http.Client, error) {
	if !InitializedHttpClient {
		Logger.Println("INFO Initializing http client")

		options := cookiejar.Options{
			PublicSuffixList: publicsuffix.List,
		}

		jar, err := cookiejar.New(&options)
		if !errorExists(err) {
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}

			HttpClient = http.Client{Transport: tr, Jar: jar}
			InitializedHttpClient = true
		} else {
			Logger.Println("ERROR Could not create HTTP Client.")
			return HttpClient, err
		}

	}
	return HttpClient, nil
}

// Function to login to wso2 carbon server using default credentials
func LoginToCarbonServer(ip string, port string, user string, pass string) bool {
	client, err := GetHttpClient()
	if err != nil {
		return false
	}
	carbonLoginUrl := "https://" + ip + ":" + port + LoginURL
	resp, err := client.PostForm(carbonLoginUrl, url.Values{
		"password": {user},
		"username": {pass},
	})
	if (!errorExists(err)) {
		defer resp.Body.Close()

		htmlcontent, err := ioutil.ReadAll(resp.Body)
		errorExists(err)

		data := string(htmlcontent)
		if resp.StatusCode == http.StatusOK && strings.Contains(data, "Signed-in as:") {
			Logger.Println("INFO WSO2 admin user Logged in")
			return true
		} else {
			Logger.Println("WARN WSO2 admin user could not login")
		}
	}
	return false
}

// Function to test if a proxy service exists in wso2 esb
func DoesProxyServiceExist(serviceName string, ip string, pass string) bool {
	client, err := GetHttpClient()
	if err != nil {
		return false
	}

	url := "https://" + ip + ":" + pass + "/services/" + serviceName + "?tryit"
	resp, err := client.Get(url)

	if errorExists(err) {
		return false
	}

	defer resp.Body.Close()

	if resp != nil {
		if resp.StatusCode == http.StatusOK {
			Logger.Printf("INFO Proxy service '" + serviceName + "' exists. Response code: %d\n", resp.StatusCode)
			return true
		} else {
			Logger.Printf("WARN Proxy service '" + serviceName + "' does not exists. Response code: %d\n", resp.StatusCode)
		}
		resp.Body.Close()
	}
	return false
}

// Function to create a new proxy service in wso2 esb
func CreateProxyService(proxyName string, targetURL string, ip string, pass string) bool {
	client, err := GetHttpClient()
	if err != nil {
		return false
	}

	args := "?formSubmitted=true&proxyName=" + proxyName + "&targetEndpointMode=url&targetURL=" + targetURL + "&publishWsdlCombo=None&availableTransportsList=https,http,local&trp__https=https&trp__http=http"
	createServiceUrl := "https://" + ip + ":" + pass + CreateProxyURL + args
	resp, err := client.Get(createServiceUrl)

	if errorExists(err) {
		return false
	}
	defer resp.Body.Close()
	Logger.Printf("INFO Proxy service creation for service '" + proxyName + "' completed. Response code: %d\n", resp.StatusCode)
	return resp.StatusCode == http.StatusOK
}

// Function to delete an existing proxy service in wso2 esb
func DeleteProxyService(proxyName string, ip string, port string) bool {
	client, err := GetHttpClient()
	if err != nil {
		return false
	}

	args := "?pageNumber=0&serviceGroups=" + proxyName
	deleteServiceUrl := "https://" + ip + ":" + port + DeleteProxyURL + args
	resp, err := client.Get(deleteServiceUrl)

	if errorExists(err) {
		return false
	}

	defer resp.Body.Close()
	Logger.Printf("INFO Proxy service deletion for service " + proxyName + " completed. Response code: %d\n", resp.StatusCode)
	return resp.StatusCode == http.StatusOK
}

func errorExists(e error) bool {
	if e != nil {
		Logger.Println("Error: " + e.Error())
		return true
	}
	return false
}
