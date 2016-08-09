package smoketests

import (
	"github.com/abhishek0198/test-framework/common"
)

func RunESBSmokeTests() {
	common.Logger.Println("Running WSO2ESB Smoke tests")
	TestProxyServiceCreationAndRemoval()
}

func TestProxyServiceCreationAndRemoval() {
	common.LoginToCarbonServer()
	common.CreateProxyService("test", "http://test.com")
	common.DoesProxyServiceExist("test")
	common.DeleteProxyService("test")
}
