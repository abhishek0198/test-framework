#!/bin/bash
set -e
DOCKERFILE_HOME=/home/abhishek/dev/dockerfiles

source common-utils.sh
source docker-utils.sh

# read script parameters
clean_previous_build=true

while getopts :n:v:r FLAG; do
    case $FLAG in
        n)
            product_name=$OPTARG
            ;;
        v)
            product_version=$OPTARG
            ;;
        r)
            clean_previous_build=false
            ;;
        \?) 
            display_usage
            ;;
    esac
done

product_path="$DOCKERFILE_HOME/$product_name"
pushd $product_path >> /dev/null

stop_docker_container

# clean existing docker build unless specified otherwise, 
# Also do a new build as well
if $clean_previous_build; then 
    clean_docker_image
    
    echo 
    echo "Starting building image..."
    echo
    bash build.sh -v $product_version
fi

echo 
echo "Starting running image..."
echo
bash run.sh -v $product_version -s

popd >> /dev/null 

echo
echo "Build and run successful"
echo

# sleep 30 seconds before running tests to allow container to complete setup
# Ports were not responding when I tried to connect immediately
echo "Waiting 30 seconds to allow container to complete setup, before running the tests"
sleep 30

# check if ports are open 
check_ports

# check if carbon server is up
check_wso2_carbon_server_status

# check wso2 carbon logs for errors
check_wso2_carbon_logs
