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
test_script_path=$(pwd)

pushd $product_path >> /dev/null

stop_docker_container

# clean existing docker build unless specified otherwise, 
# Also do a new build as well
if $clean_previous_build; then 
    clean_docker_image
    
    echo 
    echo "Starting building image..."
    echo
    bash build.sh -v $product_version > "$test_script_path/$buildlogs" 2>&1
    check_build_logs
fi

echo 
echo "Starting running image..."
echo
echo "n n" | bash run.sh -v $product_version > "$test_script_path/$runlogs" 2>&1
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
check_ports

# check if carbon server is up
check_wso2_carbon_server_status

# check wso2 carbon logs for errors
check_wso2_carbon_logs
