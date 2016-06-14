#!/bin/bash
set -e
DOCKERFILE_HOME=/home/abhishek/dev/dockerfiles
source helper.sh

# read script parameters
clean_previous_build=true

while getopts :n:v:ru FLAG; do
    case $FLAG in
        n)
            product_name=$OPTARG
            ;;
        v)
            product_version=$OPTARG
            ;;
        ru)
            clean_previous_build=false
            ;;
        \?) 
            display_usage
            ;;
    esac
done

if [[ -z ${product_name} ]] || [[ -z ${product_version} ]]; then
   display_usage
fi

product_path="$DOCKERFILE_HOME/$product_name"

if $clean_previous_build; then 
    clean_existing_images "${product_name}" "${product_version}"
fi

pushd $product_path >> /dev/null

echo 
echo "Starting building image..."
echo
bash build.sh -v $product_version 

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

