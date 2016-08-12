package smoketests

import (
	"strings"
)

func RunSmokeTest(productName string) {
	if strings.EqualFold(productName, "wso2esb") {
		RunESBSmokeTests()
	} else if strings.EqualFold(productName, "wso2mb") {
		// run message broker smoke tests
	}
}
