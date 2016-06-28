#!/bin/bash
# ------------------------------------------------------------------------
#
# Copyright 2016 WSO2, Inc. (http://wso2.com)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License

# ------------------------------------------------------------------------

set -e

###### GLOBAL VARS START ##########
DOCKERFILE_HOME=/home/abhishek/dev/dockerfiles
buildlogs="buildlogs.logs"
runlogs="runlogs.log"
###### GLOBAL VARS END ###########

TEST_START_TIME=$(date +%s)

source common-utils.sh
source docker-utils.sh

# read script parameters
clean_previous_build=true

while [[ $# -gt 1 ]]
do
key="$1"

case $key in
    -n|--name)
    product_name="$2"
    shift 
    ;;
    -v|--version)
    product_version="$2"
    shift 
    ;;
    -r|--provision-method)
    provisioning_method="$2"
    shift
    ;;
    -f|--output-file)
    test_result_file="$2"
    shift
    ;;
    *)
    display_usage
    ;;
esac
shift 
done

product_path="$DOCKERFILE_HOME/$product_name"
test_script_path=$(pwd)

echo "Running tests for $product_name, $product_version using profile $provisioning_method"

pushd $product_path >> /dev/null

echo "Stopping any running docker containers"
stop_docker_container

# clean existing docker build unless specified otherwise, 
# Also do a new build as well
 
clean_docker_image
    
echo
echo "Starting building image..."
echo
bash build.sh -v $product_version -r $provisioning_method > "$test_script_path/$buildlogs" 2>&1

echo "Checking docker build logs"
check_build_logs


echo
echo "Starting running image..."
echo
echo "n n" | bash run.sh -v $product_version > "$test_script_path/$runlogs" 2>&1

echo "Checking docker run logs"
check_run_logs

popd >> /dev/null 

echo
echo "Build and run successful"
echo

# sleep 30 seconds before running tests to allow container to complete setup
# Ports were not responding when I tried to connect immediately
echo "Waiting 30 seconds to allow container to complete setup, before running the tests"
sleep 30

# check if ports are open 
echo "Checking for exposed ports"
check_ports
echo

# check if carbon server is up
echo "Checking Carbon server status"
check_wso2_carbon_server_status
echo

# check wso2 carbon logs for errors
echo "Checking Carbon server logs for errors"
check_wso2_carbon_logs
echo

# cleanup
stop_docker_container
clean_docker_image

TEST_END_TIME=$(date +%s)
echo "Test completed in $[$TEST_END_TIME-$TEST_START_TIME] seconds"
