# WSO2 Dockerfiles test framework

WSO2 Dockerfiles provides Dockerfiles for many WSO2 products and allows users to build and run Docker images of these products. The WSO2 Dockerfiles project is still under development; changes are constantly being applied to the codebase by multiple developers of WSO2 community. In order to ensure the availability and robustness of WSO2 Dockerfiles, we are developing a general purpose integration test framework for Dockerfiles, which aims at ensuring that changes made to WSO2 Dockerfiles does not break existing functionaly and there are no regressions introduced.

## Running standard tests
Set Dockerfiles home in test-launcher.sh  
Launch test using ```test-launcher.sh -n <product-name> -v <product-version>i```  
You can specify -r if you would like to use an existing docker image for this product.

For instance to test WSO2ESB in Dockerfiles, use the following:  
```./test-launcher.sh -n wso2esb -v 4.9.0```


