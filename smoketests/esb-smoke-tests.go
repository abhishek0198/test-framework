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
	"github.com/abhishek0198/wso2dockerfiles-test-framework/common"
)

func RunESBSmokeTests() bool {
	common.Logger.Println("INFO Running WSO2ESB Smoke tests")
	result := true

	result = TestProxyServiceCreationAndRemoval() && result

	return result
}

func TestProxyServiceCreationAndRemoval() bool {
	ip := common.GetDockerContainerIP("esb")

	result := true
	result = common.LoginToCarbonServer(ip) && result
	result = common.CreateProxyService("test", "http://test.com", ip) && result
	result = common.DoesProxyServiceExist("test", ip) && result
	result = common.DeleteProxyService("test", ip) && result

	return result
}
