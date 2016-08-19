# WSO2 Dockerfiles test framework

WSO2 Dockerfiles provides Dockerfiles for many WSO2 products and allows users to build and run Docker images of these products. The WSO2 Dockerfiles project is still under development; changes are constantly being applied to the codebase by multiple developers of WSO2 community. In order to ensure the availability and robustness of WSO2 Dockerfiles, we are developing a general purpose integration test framework for Dockerfiles, which aims at ensuring that changes made to WSO2 Dockerfiles does not break existing functionaly and there are no regressions introduced.

## Build 
These instuctions are sepecific to Mac OSX, there should be similar counter part for Linux.  
+ Install Golang
  * Install using home brew ` brew install go --cross-compile-common `
  * Create default workspace folder ` mkdir $HOME/go `
  * Setup paths in your .bash_profile  
     ` export GOPATH=$HOME/go `  
     ` export PATH=$PATH:$GOPATH/bin ` 
+ Get the project source
   ` go get -d github.com/abhishek0198/wso2dockerfiles-test-framework `
+ Build the project,  
  * cd $GOPATH/src/github.com/abhishek0198/test-framework  
  * go build
+ Lauch the test framework using  
  ` ./wso2dockerfiles-test-framework `

## Running standard tests
The test framework also requires setting up project relevent to your tests. Following are the projects that you should clone:  
WSO2 Dockerfiles (https://github.com/wso2/dockerfiles)  
WSO2 Puppet Modules (https://github.com/wso2/puppet-modules)  

You will also need to download java and product specific zip files. Instructions can be found on WSO2 Dockerfiles.  

Once above setup is completed, follow following steps to run the tests:  
1. Edit TestConfigPath under <project_root>/src/common/test-config.json  
2. Set dockerfileshome and carob_server_port in <project_root>/src/config/test-config.json  
3. Configure products to test along with desired provisioning in <project_root>/src/config/test-config.json  
4. Launch test using ```./wso2dockerfiles-test-framework``` from the same directory  

## Mac and Windows users
Docker does not run directly on Mac OSX and Windows, instead you will need to use docker-machine tool. From Docker's official page:  

"Docker Machine is a tool that lets you install Docker Engine on virtual hosts, and manage the hosts with docker-machine commands. You can use Machine to create Docker hosts on your local Mac or Windows box, on your company network, in your data center, or on cloud providers like AWS or Digital Ocean."

Follow the follwing instructions to setup test framework on Mac or Windows:  

1. Install docker-machine using:  
https://docs.docker.com/machine/install-machine/  

2. Create default docker host  
`` docker-machine create --driver virtualbox default ``  

3. Now, start a docker host  
`` docker-machine start default ``  

4. The command complets with  
`` Started machines may have new IP addresses. You may need to re-run the 'docker-machine env' command. ``  
I have this command added to ~/.bash_profile so that when I start my shell, and docker host is already created, I do not need to run the env command.  
`` eval $(docker-machine env default) ``  

5. You will need the IP address of this newly created docker host and configure it on the test-config. You can do it using `` docker-machine inspect default `` and get the IPAddress part.  

## Unix/Linux users  
In order to support running on Mac OSX, `carbon_server_ip` is explicitly set to use docker-machine created host IP. Remove this config, if you're running on Linux

## Sample test config  
```        
{
   "testconfig":{
      "wso2_products":[
         {
            "enabled":"true",
            "name":"wso2esb",
            "version":"4.9.0",
            "provisioning_method":["default","puppet"],
            "platform":"default"
         },
         {
            "enabled":"true",
            "name":"wso2mb",
            "version":"3.1.0",
            "provisioning_method":["puppet"],
            "platform":"default"
         },
         {
            "enabled":"false",
            "name":"wso2esb",
            "version":"4.9.0",
            "provisioning_method":["puppet"],
            "platform":"kubernetes",
            "profile":["worker", "manager"]
         },
         {
            "enabled":"false",
            "name":"wso2am",
            "version":"2.0.0",
            "provisioning_method":["default"],
            "platform":"default"
         }
      ],
      "output_file":"dockerfiles-test-result.txt",
      "dockerfileshome":"/Users/abhishektiwari/dev/dockerfiles",
      "carbon_server_ip":"192.168.99.100",
      "carbon_server_port":"9443",
      "carbon_server_username":"admin",
      "carbon_server_password":"admin"
   }
}
```
The config above is to test WSO2ESB using default and puppet provisioning methods. It also tests WSO2MB using puppet provisioning.
