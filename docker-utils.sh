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

# get the IP address of the running docker container for the product under test
function get_docker_container_ip() {
    container_id=$(docker ps | grep $product_name | awk '{print $1}')
    container_ip=$(docker inspect --format '{{ .NetworkSettings.IPAddress }}' "${container_id}")
    if [ -z "${container_ip}" ]; then
        echo "Could not find IP address of container ${container_id}"
        exit 1
    fi
    echo $container_ip
}

# clean a docker image using product tag
function clean_docker_image() {
   product_tag="$product_name:$product_version"
   if docker history -q $product_tag > /dev/null 2>&1; then
       echo "Removing docker image $product_tag"
       docker rmi $product_tag >> /dev/null
   fi
}

# stop a running docker container using product tag before running tests.
function stop_docker_container() {
   container_id=$(docker ps -a | grep $product_name | awk '{print $1}')
   # if container exists, stop and remove the container
   if [ -n "$container_id" ]; then
       echo "Stoping container $container_id" 
       docker rm $(docker stop $container_id) >> /dev/null
   fi
}

function copy_carbon_logs() {
    container_id=$(docker ps | grep $product_name | awk '{print $1}')
    ip=$(get_docker_container_ip)
    docker cp "$container_id:/mnt/$ip/$product_name-$product_version/repository/logs/" ./
}

