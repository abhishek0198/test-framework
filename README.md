# WSO2 Dockerfiles test framework

WSO2 Dockerfiles provides Dockerfiles for many WSO2 products and allows users to build and run Docker images of these products. The WSO2 Dockerfiles project is still under development; changes are constantly being applied to the codebase by multiple developers of WSO2 community. In order to ensure the availability and robustness of WSO2 Dockerfiles, we are developing a general purpose integration test framework for Dockerfiles, which aims at ensuring that changes made to WSO2 Dockerfiles does not break existing functionaly and there are no regressions introduced.

## Running standard tests
Set Dockerfiles home in test-launcher.sh  
Configure products to test along with desired provisioning in test-config.json
Launch test using ```python test-runner.py```  

For instance to test WSO2ESB using puppet provisioning, use the following:  
```        
	"products": {
                "esb": {
                        "enabled":"true",
                        "name": "wso2esb",
                        "version": "4.9.0",
                        "provisioning_method": "puppet"
                }
        }
```


