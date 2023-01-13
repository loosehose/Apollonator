# Apollonator
## Description

Apollonator is a script that extracts information from Apollo.io for target organizations using a list of names. It utilizes a configuration file (config.yml) to input the API key, organization name, and boolean values for specific information to gather from the Apollo JSON response. To avoid usage limitations, Apollonator implements a 18 second delay between each request.

### API Information
I created a master key to avoid any issues; however, that may be overkill.

Apollo's API key can be found here: https://developer.apollo.io/keys/

### Formatting

The name file should be in a JOHN SMITH format. Otherwise, the script will break. 

### Disclaimer

During the production of this script, I used the 'Professional' subscription becauses I had a large list of users. Depending on your volume, you may be able to get away with the 'Basic' plan. This script is also not bug-proof. This was created for a red team engagement and quickly thrown together.

## Quick Start

### Installation

```
git clone https://github.com/loosehose/Apollonator.git
cd Apollonator
pip3 install -r requirements
```

### Basic 

```
python3 apollonator.py -c config.yml -n users.txt
```

### Create Excel Sheet

```
python3 apollonator.py -c config.yml -n users.txt -e
```

