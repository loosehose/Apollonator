import requests
import json
from json.decoder import JSONDecodeError
import time
import argparse
import os
import openpyxl
from openpyxl.utils.dataframe import dataframe_to_rows
import yaml
import pandas as pd

class Apollonator():
    def validate_args(self):
        parser = argparse.ArgumentParser(description="")
        parser.add_argument("-c", "--config", dest="config", default=None, 
                            help="Import the config.yml file with updated information.")
        parser.add_argument("-e", "--excel", dest="excel", action='store_true',
                            help="Save the results to an excel file.")
        parser.add_argument("-n", "--names", dest="name", default=None,
                            help="Input a list of names.")
        parser.add_argument("-s", "--sleep", dest="sleep", default=18, type=int,
                        help="Specify sleep delay in seconds. Default is 18.")
        args = parser.parse_args()

        if not (args.config):
            print("[-] Specify an input file [config.yml] with -c")
            exit(1)
        if args.config != None:
            if os.path.exists(args.config) == False:
                print("[-] " + args.config + " does not exist in directory")
                exit(1)
        return args
    
    def yaml_parser(self, config):
        with open(config, "r") as stream:
            config = yaml.safe_load(stream)["apollonator"]
            return [config["api_key"], config["organization"], config["email"], config["title"]]

    def apollo_requester(self, api_key, organization_name, first_name, last_name, delay):
        url = "https://api.apollo.io/v1/people/match"
        payload = { 
            "api_key": api_key, 
            "first_name": first_name, 
            "last_name": last_name, 
            "organization_name": organization_name, 
        }

        r = requests.post(url, json=payload)
        time.sleep(delay)
        return r.text

    def extract_email_from_json(self, json_file):
        email = json.loads(json_file)

        return email["person"]["email"]

    def extract_title_from_json(self, json_file):
        title = json.loads(json_file)
        return title["person"]["title"]
                
    def convert_to_excel(self, first_name, last_name, org, email, title):
        file = "appolonator"+org+".xlsx"
        if email != None:
            domain = email[email.index('@') + 1 : ]
        else:
            email = ""
            domain = ""

        apolloDf = pd.DataFrame()
        apolloDf = pd.concat([apolloDf, pd.DataFrame.from_records([{'First Name': first_name, "Last Name": last_name, "Organization": org, 'Email': email, 'Domain': domain, 'Title': title}])])

        if os.path.isfile(file):  # if file already exists append to existing file
            workbook = openpyxl.load_workbook(file)  
            sheet = workbook['Employee Info']  

            # append the dataframe results to the current excel file
            for row in dataframe_to_rows(apolloDf, header = False, index = False):
                sheet.append(row)
            workbook.save(file)  
            workbook.close()  
        else:  # create the excel file if doesn't already exist
            with pd.ExcelWriter(path = file, engine = 'openpyxl') as writer:
                apolloDf.to_excel(writer, index = False, sheet_name = 'Employee Info')

    def run(self, args):
        config_check = self.yaml_parser(args.config)
        api_key, org, boolEmail, boolTitle  = config_check[0], config_check[1], config_check[2], config_check[3]

        with open(args.name) as f:
            lines = f.readlines()
            est_time = len(lines) * 18 / 60  # estimated completion time in minutes
            print(f"Estimated completion time: {est_time} minutes.")
            for name in lines:
                first_name = name.split()[0]
                last_name = name.split()[1]

                # Skip LinkedIn Member
                if first_name.lower() == 'linkedin' and last_name.lower() == 'member':
                    continue

                requester = self.apollo_requester(api_key, org, first_name, last_name, args.sleep)
                try:
                    if boolEmail == True: 
                        email = self.extract_email_from_json(requester)
                        if email is None or email.lower() == 'none':
                            continue

                    if boolTitle == True: 
                        title = self.extract_title_from_json(requester)
                    
                    print(first_name, last_name, email)
                    if args.excel:
                        Apollonator().convert_to_excel(first_name, last_name, org, email, title)
                    
                except JSONDecodeError as e:
                    print("API Daily Limit Reached")
                    exit(1)
            

if __name__ == "__main__":
    a = Apollonator()
    args = a.validate_args()
    a.run(args)
