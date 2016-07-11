package common

const (	
	BuildLogs = "buildlogs.log"
	RunLogs = "runlogs.log"
	TestConfigPath = "/home/abhishek/dev/test-framework/src/config/test-config.json"
)

var (
	Testconfig TestConfig
	DockerContainerID string
	DockerContainerIP string 
)

type Product struct {
	Enabled   string
    Name   string
    Version   string
    Provisioning_method   string
    Organization string
    Platform string
}

type TestConfig struct {
	Wso2_products []Product
    Output_file string
    DockerfilesHome string
	Carbon_Server_Port string
}

func ResetTestSpecificVariables() {
	DockerContainerID = ""
	DockerContainerIP = ""
}