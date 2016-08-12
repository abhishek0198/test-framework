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

package smoketests

import (
	"github.com/abhishek0198/test-framework/common"
)

func RunESBSmokeTests() {
	common.Logger.Println("Running WSO2ESB Smoke tests")
	TestProxyServiceCreationAndRemoval()
}

func TestProxyServiceCreationAndRemoval() {
	ip := common.GetDockerContainerIP("esb")

	common.LoginToCarbonServer(ip)
	common.CreateProxyService("test", "http://test.com", ip)
	common.DoesProxyServiceExist("test", ip)
	common.DeleteProxyService("test", ip)
}
