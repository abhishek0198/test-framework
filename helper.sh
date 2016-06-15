#!/bin/bash

# show usage of test framework
function display_usage() {
   echo "Insufficient or invalid options provided!"
   echo 
   echo "Options:"
   echo -en " -n\t"
   echo "[REQUIRED] Product name to test"
   echo -en " -v\t"
   echo "[REQUIRED] Product version to test"
   echo -en " -r\t"
   echo "[OPTIONAL] Reuse existing product build"
   echo
   exit 1
}

# clean a docker image using product tag
function clean_docker_image() {
   product_tag="${1}:${2}"
   if docker history -q $product_tag > /dev/null 2>&1; then
       echo "Removing docker image $product_tag"
       docker rmi $product_tag >> /dev/null
   fi
}

# stop a running docker container using product tag before running tests.
function stop_docker_container() {
   product_tag="${1}:${2}"
   container_id=$(docker ps -a -q --filter ancestor="$product_tag")
   
   # if container exists, stop and remove the container
   if [ -n "$container_id" ]; then
       echo "Stoping container $container_id" 
       docker rm $(docker stop $container_id) >> /dev/null 
   fi
}

# check all the ports exposed in the Dockerfile
# TODO: read them from dockerfile directly?
function check_ports() {
    declare -a ports=("8280" "8243" "9763" "9443")
    host="172.17.0.2"
    for port in "${ports[@]}"
    do 
        echo "checking $port"
        port_closed=$(curl -s $host:$port > /dev/null && echo false || echo true)

        if $port_closed; then
            echo "Unable to connect to  $host:$port."
            #exit 1
        else 
            echo "Connection to $host:$port successful."
        fi
     done
}

# check if the web server has been started successfully
# TODO: read server address and port from config
function check_carbon_server() {
    http_response_code=$(curl --insecure --write-out %{http_code} --silent --output /dev/null "https://172.17.0.2:9443/carbon/admin/login.jsp")

    if [ "$http_response_code" == "200" ]; then
        echo "Carbon server is up and running."
    else
        echo "Carbon server is not running."
    fi    
}
