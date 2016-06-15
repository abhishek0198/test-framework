#!/bin/bash
set -e
DOCKERFILE_HOME=/home/abhishek/dev/dockerfiles
source helper.sh

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

stop_docker_container "${product_name}" "${product_version}"

# clean existing docker build unless specified otherwise, 
# Also do a new build as well
if $clean_previous_build; then 
    clean_docker_image "${product_name}" "${product_version}"
    
    echo 
    echo "Starting building image..."
    echo
    bash build.sh -v $product_version
fi

echo 
echo "Starting running image..."
echo
bash run.sh -v $product_version 

popd >> /dev/null 

echo
echo "Build and run successful"
echo

# check if ports are open 
check_ports

# check if carbon server is up
check_carbon_server

