#/bin/bash
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

# check all the ports exposed in the Dockerfile
function check_ports() {
    host_ip=$(get_docker_container_ip)
    port_line=$(grep "EXPOSE" "$product_path/Dockerfile")
    read -a port_array <<<$port_line
    ports_to_check=("${port_array[@]:1}")
    for port in "${ports_to_check[@]}"
    do
        nc -z -v -w5 $host_ip $port >> $test_result_file 2>&1
    done
    echo 
}

# check if the web server has been started successfully
function check_wso2_carbon_server_status() {
    server_port=9443
    server_host=$(get_docker_container_ip)
    http_response_code=$(curl --insecure --write-out %{http_code} --silent --output /dev/null "https://$server_host:$server_port/carbon/admin/login.jsp")

    if [ "$http_response_code" == "200" ]; then
        echo "Carbon server is up and running."
    else
        echo "Carbon server is not running."
    fi   
}

# check wso2 carbon logs from the running container for any errors
function check_wso2_carbon_logs() {
    # copy logs from Docker container to local
    copy_carbon_logs
    #pushd logs >> /dev/null
    #errors=$(grep -ir 'error' .)
    #popd >> /dev/null
    echo
    if $(grep -ri 'error' "./logs"); then
        echo "WSO2 Carbon logs contain errors. Please verify them in ./logs/."
    else  
        echo "WSO2 Carbon logs does not contain any errors."
        rm -rf "./logs"
    fi
    echo 
}

function check_build_logs() {
    if $(grep -i 'error' "$test_script_path/$buildlogs"); then
        echo "Docker build has errors. Please verify them in $buildlogs"
    else
        echo "No errors found in $buildlogs"
        rm "$test_script_path/$buildlogs"
    fi
}

function check_run_logs() {
    if $(grep -i 'error' "$test_script_path/$runlogs"); then    
        echo "Docker container run has errors. Please verify them in $runlogs"
    else
        echo "No errors found in $runlogs"
        rm "$test_script_path/$runlogs"
    fi    
}

