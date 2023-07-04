<p align="center">
  <img src="https://user-images.githubusercontent.com/75705022/212420146-b2ccb43b-f803-49a9-a362-50ba4e789048.png" />
</p>

Apollonator is a Golang utility that sends requests to the Apollo.io API to retrieve person's email and title. It reads the person's name and organization from a text file and returns the extracted email and title. The purpose of this tool is to enumerate target email addresses without having to touch their infrastructure. 

The retrieved information can optionally be saved into an Excel file for further use.

## Installation

```
git clone https://github.com/loosehose/apollonator.git
cd Apollonator
go build
```

## Configuration

Before running the program, you need to configure your [apollo.io API key](https://developer.apollo.io/keys/) and the organization name in the config.yml file. The organization name can quickly be found by searching for it using the Apollo website searching engine.

*Note: sometimes you really need to use that searching mechanic on Apollo because the organization name may be slightly different than anticipated.*

```
apollonator:
  api_key: "YOUR_API_KEY"
  organization: "YOUR_ORGANIZATION_NAME"
  email: true
  title: true
```

Replace "YOUR_API_KEY" with your Apollo.io API key and "YOUR_ORGANIZATION_NAME" with your organization name. Set the email and title fields to true or false depending on whether you want to extract email or title respectively. Apollo can also give you phone numbers. Originally, I added this; however, it usually just returns the company phone number... which is kinda useless.

## Usage
```
./apollonator -h                                                                       

        Holy OSINT, Apollonator! 
        author: github.com/loosehose

Usage of ./apollonator:
  -c string
        Import the config.yml file with updated information.
  -e    Save the results to an excel file.
  -n string
        Input a list of names.
  -s int
        Specify sleep delay in seconds. (default 18)
```

You can run the program with the following command:

```
./apollonator -n /path/to/users.txt -c ./config -e
```

The `-c` or `--config` option specifies the configuration file.
The `-n` or `--names` option specifies the text file that contains the list of names to be looked up. Each line of the file should contain a first and last name separated by a space.
If you want the results saved in an Excel file, add the `-e` or `--excel option`.
The `-s` or `--sleep` option specifies the delay between each request. 18 is the default.

## Output

The program prints out the first name, last name, and email (if requested) of each person in the provided names file. If the `-e` option is used, the information is saved in an Excel file named "apollonator_{organization}.xlsx", where "{organization}" is replaced with the name of your organization. The Excel file has a sheet named "Employee Info" with columns for First Name, Last Name, Organization, Email, Domain (extracted from the email), and Title (if requested). If the Excel file already exists, the new information is appended to it. One neat thing about the excel file is it will automatically pick the domain name out. This makes it easier to sort out different domain names the company may be using and filter the false positives out.

If you want to copy the output without the logging, simply change line 249 (and rebuild):

```go
// current
log.Info().Msgf("%s %s %s", firstName, lastName, response.Person.Email)
```

```go
// updated
fmt.Printf("%s %s %s", firstName, lastName, response.Person.Email)
```

## Caveats

Let's say you have a user named Jane Smith Johnson. While this name may be within the Apollo database, it is wise to create three names from this: Jane Smith Johnson, Jane Smith, and Jane Johnson to ensure you get a match. This edge case, as well as a few others, will be fixed in future updates, *hopefully*. For now, Apollonator will look for Jane Johnson only in this case.

## Error Handling

- If the program reaches your daily API limit it will give an ERR: API limit reached. Please try again later.
