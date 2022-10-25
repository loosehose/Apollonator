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
            return [config["api_key"], config["organization"], config["email"], config["phone_number"], config["title"]]

    def apollo_requester(self, api_key, organization_name, first_name, last_name):
        url = "https://api.apollo.io/v1/people/match"
        payload = { "api_key": api_key, "first_name": first_name, "last_name": last_name, "organization_name": organization_name, }

        r = requests.post(url, json=payload)
        time.sleep(18)
        return r.text

    def extract_email_from_json(self, json_file):
        email = json.loads(json_file)

        return email["person"]["email"]

    def extract_phone_from_json(self, json_file):
        json_loader = json.loads(json_file)

        try:
            tele = json_loader["person"]['phone_numbers'][0]
            return tele["sanitized_number"]
        except KeyError:
            return

    def extract_title_from_json(self, json_file):
        title = json.loads(json_file)
        return title["person"]["title"]

    def get_names_and_store_in_txt(self, names):
        # reset fullnames.txt
        file = time.strftime("%Y%m%d%H%M")+"_names.txt"
        fullnames = open(file, "w")

        with open(names) as f:
            lines = f.readlines()
            for name in lines:
                name = name.replace(',', ' ')
                first_name = name.split()[1]
                last_name = name.split()[0]
                f_name = first_name + " " + last_name + "\n"
                fullnames = open(file, "a")
                fullnames.write(f_name)
    def convert_to_excel(self, first_name, last_name, org, email, phone, title):
        file = "appolonator"+org+".xlsx"
        if email != None:
            domain = email[email.index('@') + 1 : ]
        else:
            email = ""
            domain = ""

        apolloDf = pd.DataFrame()
        apolloDf = pd.concat([apolloDf, pd.DataFrame.from_records([{'First Name': first_name, "Last Name": last_name, "Organization": org, 'Email': email, 'Domain': domain, 'Phone Number': phone, 'Title': title}])])

        if os.path.isfile(file):  # if file already exists append to existing file
            workbook = openpyxl.load_workbook(file)  # load workbook if already exists
            sheet = workbook['Employee Info']  # declare the active sheet 

            # append the dataframe results to the current excel file
            for row in dataframe_to_rows(apolloDf, header = False, index = False):
                sheet.append(row)
            workbook.save(file)  # save workbook
            workbook.close()  # close workbook
        else:  # create the excel file if doesn't already exist
            with pd.ExcelWriter(path = file, engine = 'openpyxl') as writer:
                apolloDf.to_excel(writer, index = False, sheet_name = 'Employee Info')

    def run(self, args):
        config_check = Apollonator().yaml_parser(args.config)
        api_key, org, boolEmail, boolPhone, boolTitle = config_check[0], config_check[1], config_check[2], config_check[3], config_check[4]

        Apollonator().get_names_and_store_in_txt(args.name)
        file = time.strftime("%Y%m%d%H%M")+"_names.txt"
        with open(file) as f:
            lines = f.readlines()
            for name in lines:
                first_name = name.split()[0]
                last_name = name.split()[1]
                requester = Apollonator().apollo_requester(api_key, org, first_name, last_name)
                try:
                    if boolEmail == True: email = Apollonator().extract_email_from_json(requester)
                    if boolPhone == True: phone = Apollonator().extract_phone_from_json(requester)
                    if boolTitle == True: title = Apollonator().extract_title_from_json(requester)
                    print(first_name, last_name, org, email, phone, title)
                    if args.excel:
                        Apollonator().convert_to_excel(first_name, last_name, org, email, phone, title)
                    
                except JSONDecodeError as e:
                    print("API Daily Limit Reached")
                    exit(1)
                
                

if __name__ == "__main__":
    a = Apollonator()
    args = a.validate_args()
    a.run(args)