<p align="center">
  <img src="https://user-images.githubusercontent.com/75705022/250726305-e20cb22e-9822-42e0-bb37-b7e73651b2be.png" />
</p>

Apollonator is a Golang utility that uses the Apollo.io API to retrieve a person's email and job title. By reading in a list of names and an organization from a text file, it is able to return corresponding emails and job titles. This utility enables the enumeration of target email addresses without interacting with the target's infrastructure.

An optional feature allows the retrieved information to be saved into an Excel file for later use.

## Installation
Clone the repository, navigate to the 'Apollonator' directory, and build the program:
```bash
git clone https://github.com/loosehose/apollonator.git
cd Apollonator
go build
```

## Configuration
Before running the program, acquire your [apollo.io API key](https://developer.apollo.io/keys/) and specify your organization name in the config.yml file. You can create as many accounts as you want for free and add multiple API keys in the config.yml file. Apollonator will adjust the default sleep based on this.

```yaml
apollonator:
  api_keys:
    - "FirstAPIKey"
    - "SecondAPIKey"
    - "ThirdAPIKey"
    # Add more keys as needed
  organization: "Github"
  email: true
  title: true
```
*Note: sometimes you really need to use that searching mechanic on Apollo because the organization name may be slightly different than anticipated.*

Replace "YOUR_API_KEY" with your Apollo.io API key and "YOUR_ORGANIZATION_NAME" with your organization's exact name. Set the email and title fields to true or false based on your requirements.

## Usage
```
./apollonator
                  
                   =@@#                 
                  +@@@@%.               
                 *@@@@@@@:   .....      
                #@@@@@@@@@-  :---       
              .%@@@@@@@@@@@=  ::        
             :@@@@@@-.%@@@@@*           
            -@@@@@@:   *@@@@@#          
           =@@@@@%.  -  +@@@@@#         
          *@@@@@#  .#@*  =@@@@@%:       
         #@@@@@*  .%@@@%. :@@@@@@-      
       .%@@@@@=  :@@@@@@:  .%@@@@@=     
      .%@@@@@-  -@@@@@@:     #@@@@@+    
     -@@@@@@:  =@@@@@%.       #@@@@@#   
    =@@@@@%.  +@@@@@#          *@@@@@%. 
   +@@@@@@+=*@@@@@@*            =@@@@@%.
  *@@@@@@@@@@@@@@@+              -@@@@@%
 #@@@@@@@@@@@@@%+.                :@@@@% (pollonator)


        author: github.com/loosehose

[-c|--config] is required
usage: apollonator [-h|--help] -c|--config "<value>" [-e|--excel] [-n|--names
                   "<value>"] [-s|--sleep <integer>]

                   

Arguments:

  -h  --help    Print help information
  -c  --config  Import the config.yml file with updated information.
  -e  --excel   Save the results to an excel file.
  -n  --names   Input a list of names.
  -s  --sleep   Specify sleep delay in seconds.. Default: 18
```

You can run the program with the following command:

```
./apollonator -n /path/to/users.txt -c ./config.yml -e
```

The `-c` or `--config` option specifies the configuration file.
The `-n` or `--names` option specifies the text file that contains the list of names to be looked up. Each line of the file should contain a first and last name separated by a space.
If you want the results saved in an Excel file, add the `-e` or `--excel option`.
The `-s` or `--sleep` option specifies the delay between each request. 18 is the default. You probably shouldn't change this number or Apollo's API will put you in time out. If you use multiple API keys from different accounts (wink wink), the default sleep will be adjusted. For instance, if you have two API keys, the sleep will default to 9 seconds.

## Output

The utility prints the first name, last name, and requested email of each person in the provided names file. If the -e option is used, this data is saved into an Excel file named "apollonator_{organization}.xlsx", where "{organization}" is replaced with the name of your organization. There is regex in place for edge cases for organizations such as O'Reilly.

The Excel file contains a sheet named "Employee Info" with columns for First Name, Last Name, Organization, Email, Domain (extracted from the email), and Title (if requested). If the file already exists, new data is appended.

If you want to copy the output without the logging, simply change the following line (and rebuild):

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
