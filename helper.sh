#/bin/bash

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
        nc -z -v -w5 $host_ip $port
    done 
}

# check if the web server has been started successfully
function check_carbon_server() {
    server_port=9443
    server_host=$(get_docker_container_ip)
    http_response_code=$(curl --insecure --write-out %{http_code} --silent --output /dev/null "https://$server_host:$server_port/carbon/admin/login.jsp")

    if [ "$http_response_code" == "200" ]; then
        echo "Carbon server is up and running."
    else
        echo "Carbon server is not running."
    fi    
}
