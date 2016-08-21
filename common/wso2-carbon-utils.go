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
	"net/http"
	"net/http/cookiejar"
)

var (
	HttpClient            http.Client
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
			Logger.Println(GetRedColorFormattedOutputString("ERROR Could not create HTTP Client."))
			return HttpClient, err
		}

	}
	return HttpClient, nil
}

func errorExists(e error) bool {
	if e != nil {
		Logger.Println(GetRedColorFormattedOutputString("Error: " + e.Error()))
		return true
	}
	return false
}
