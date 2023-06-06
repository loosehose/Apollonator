<p align="center">
  <img src="https://user-images.githubusercontent.com/75705022/212420146-b2ccb43b-f803-49a9-a362-50ba4e789048.png" />
</p>

Apollonator is a Python utility that sends requests to the Apollo.io API to retrieve person's email and title. It reads the person's name and organization from a text file and returns the extracted email and title.

The retrieved information can optionally be saved into an Excel file for further use.

## Prerequisites
```
pip install -r requirements.txt
```
## Configuration
Before running the program, you need to configure your Apollo.io API key and the organization name in the config.yml file.
```
apollonator:
  api_key: "YOUR_API_KEY"
  organization: "YOUR_ORGANIZATION_NAME"
  email: true
  title: true
```
Replace "YOUR_API_KEY" with your Apollo.io API key and "YOUR_ORGANIZATION_NAME" with your organization name. Set the email and title fields to true or false depending on whether you want to extract email or title respectively.

## Usage
You can run the program with the following command:
```
python apollonator.py -c config.yml -n names.txt
```
The `-c` or `--config` option specifies the configuration file.
The `-n` or `--names` option specifies the text file that contains the list of names to be looked up. Each line of the file should contain a first and last name separated by a space.
If you want the results saved in an Excel file, add the `-e` or `--excel option`.
Each request to the Apollo.io API takes about 18 seconds to complete, so the estimated completion time for all names in the file is output at the start of the run.

Output
The program prints out the first name, last name, and email (if requested) of each person in the provided names file. If the `-e` option is used, the information is saved in an Excel file named "apollonator{organization}.xlsx", where "{organization}" is replaced with the name of your organization. The Excel file has a sheet named "Employee Info" with columns for First Name, Last Name, Organization, Email, Domain (extracted from the email), and Title (if requested). If the Excel file already exists, the new information is appended to it.

## Error Handling
If the program encounters a JSONDecodeError, it assumes that the daily API limit has been reached and exits.
