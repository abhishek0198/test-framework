# WSO2 Dockerfiles test framework

WSO2 Dockerfiles provides Dockerfiles for many WSO2 products and allows users to build and run Docker images of these products. The WSO2 Dockerfiles project is still under development; changes are constantly being applied to the codebase by multiple developers of WSO2 community. In order to ensure the availability and robustness of WSO2 Dockerfiles, we are developing a general purpose integration test framework for Dockerfiles, which aims at ensuring that changes made to WSO2 Dockerfiles does not break existing functionaly and there are no regressions introduced.

## Build using Eclipse
These instuctions are sepecific to Mac OSX, there should be similar counter part for Linux
1. Install Golang  
   a) Install using home brew ``` brew install go --cross-compile-common ```  
   b) Create default workspace folder ```mkdir $HOME/go```  
   c) Setup paths in your .bash_profile
     ```
     export GOPATH=$HOME/go
     export PATH=$PATH:$GOPATH/bin
     ```
2. Download and install eclipse neon or later
3. Install Goclipse plugin from update site ```http://goclipse.github.io/releases/```.  
Detailed instuctions (https://github.com/GoClipse/goclipse/blob/latest/documentation/Installation.md#installation)
4. Clone this git repo
5. Create a new Go eclipse project and use the source location of the cloned git project
6. Configure Go preferences. For Go installation use /usr/local and Eclipse go path use $HOME/go directory
7. Congigure Tools under Go Preferences, just hit Download button for all the three and it will setup those tools.
8. You can now build and run this project as a Go Application. 
   
## Build without Eclipse
1. Follow (1), and (4) from above
2. Add the project location to GOPATH
   Mine looks like  
``` /Users/abhishektiwari/go:/Users/abhishektiwari/dev/test-framework ```  
3. Build project by launching following command from project root  
  ``` go install -v -gcflags "-N -l" ./... ```  
   The built application will be in $PROJECT_ROOT/bin.
   You can launch ```./main``` from bin directory

## Running standard tests
The test framework also requires setting up project relevent to your tests. Following are the projects that you should clone:  
WSO2 Dockerfiles (https://github.com/wso2/dockerfiles)  
WSO2 Puppet Modules (https://github.com/wso2/puppet-modules)  

You will also need to download java and product specific zip files. Instructions can be found on WSO2 Dockerfiles.  

Once above setup is completed, follow following steps to run the tests:  
1. Set dockerfileshome and carob_server_port in src/config/test-config.json  
2. Configure products to test along with desired provisioning in src/config/test-config.json
3. Launch test using ```./main``` from bin directory  

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
