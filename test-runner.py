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
	
