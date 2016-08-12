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
