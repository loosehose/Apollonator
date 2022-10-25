# Apollonator

## Description

Apollonator is a script designed to extract information from Apollo.io on target organizations, given a list of names. Apollonator relies on a configuration file (config.yml) to parse the API key, organization name, and boolean values to determine what information will be gathered from the Apollo JSON response. 

Apollo's API key can be found here: https://developer.apollo.io/keys/
I created a master key to avoid any issues; however, that may be overkill.

The name file should be in a JOHN SMITH format. Otherwise, the script will break. 

### Disclaimer

During the production of this script, I used the 'Professional' subscription becauses I had a large list of users. Depending on your volume, you may be able to get away with the 'Basic' plan. 

## Quick Start

### Installation

```
git clone https://github.com/loosehose/Apollonator.git
cd Apollonator
pip3 install -r requirements
```

### Basic 

```
python3 apollonator -c config.yml -n users.txt
```

### Create Excel Sheet

```
python3 apollonator -c config.yml -n users.txt -e
```

