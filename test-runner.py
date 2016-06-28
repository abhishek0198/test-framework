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


import json
import subprocess
import shlex

with open('test-config.json') as data_file:    
    data = json.load(data_file)

# read test result output file path and open the file
filepath = data['test_result_file']

products_to_test = data['products']

for product in products_to_test.keys():
	product_settings = products_to_test[product]
	enabled = product_settings['enabled']
	if enabled == 'true':
        	product_name = product_settings['name']
		product_version = product_settings['version']
		provisioning_method = product_settings['provisioning_method']
		test_command = './test-launcher.sh -n ' + product_name + ' -v ' + product_version + ' -r ' + provisioning_method + ' -f ' + filepath
		with open(filepath, "a") as output_file:
			subprocess.call(shlex.split(test_command), stdout=output_file)
	
