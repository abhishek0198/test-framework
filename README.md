# WSO2 Dockerfiles test framework

WSO2 Dockerfiles provides Dockerfiles for many WSO2 products and allows users to build and run Docker images of these products. The WSO2 Dockerfiles project is still under development; changes are constantly being applied to the codebase by multiple developers of WSO2 community. In order to ensure the availability and robustness of WSO2 Dockerfiles, we are developing a general purpose integration test framework for Dockerfiles, which aims at ensuring that changes made to WSO2 Dockerfiles does not break existing functionaly and there are no regressions introduced.

## Running standard tests
You will need Golang install in your dev machine.
Set dockerfileshome and carob_server_port in src/config/test-config.json  
Configure products to test along with desired provisioning in src/config/test-config.json
Build the project using ```go build```
Launch test using ```./main``` from bin directory  

For instance, here is a sample test config to test WSO2ESB using default and WSO2MB using puppet provisioning.

```        
{
   "testconfig":{
      "wso2_products":[
         {
            "enabled":"True",
            "name":"wso2esb",
            "version":"4.9.0",
            "provisioning_method":"default"
         },
         {
            "enabled":"False",
            "name":"wso2mb",
            "version":"3.1.0",
            "provisioning_method":"puppet",
            "platform":"default"
         }
      ],
      "output_file":"/home/abhishek/dev/test-framework/output.txt",
      "dockerfileshome":"/home/abhishek/dev/dockerfiles",
      "carbon_server_port":"9443"
   }
}
```


