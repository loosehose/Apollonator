package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/akamensky/argparse"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tealeg/xlsx"
	"gopkg.in/yaml.v2"
)

const (
	apolloURL         = "https://api.apollo.io/v1/people/match"
	defaultSleepDelay = 18
	outputExcelFile   = "apollonator.xlsx"
)

var apiKeyIndex int

type Config struct {
	APIKeys      []string `yaml:"api_keys"`
	Organization string   `yaml:"organization"`
	Email        bool     `yaml:"email"`
	Title        bool     `yaml:"title"`
}

type ApolloRequest struct {
	APIKey           string `json:"api_key"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	OrganizationName string `json:"organization_name"`
}

type ApolloResponse struct {
	Person struct {
		Email string `json:"email"`
		Title string `json:"title"`
	} `json:"person"`
}

type PersonData struct {
	FirstName    string
	LastName     string
	Organization string
	Email        string
	Domain       string
	Title        string
}

func ParseYaml(configFile string) (Config, error) {
	config := struct {
		Apollonator Config `yaml:"apollonator"`
	}{}

	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return Config{}, err
	}

	if len(config.Apollonator.APIKeys) == 0 {
		return Config{}, errors.New("missing API keys in the configuration")
	}

	if config.Apollonator.Organization == "" {
		return Config{}, errors.New("missing Organization in the configuration")
	}

	return config.Apollonator, nil
}

func ApolloRequester(apollonator *Config, firstName string, lastName string, delay time.Duration, apiKeyIndex int) (ApolloResponse, error) {
	// Check if there are any API keys left
	if len(apollonator.APIKeys) == 0 {
		return ApolloResponse{}, errors.New("All API keys are out of gas! Exiting.")
	}

	// Rotate the API key index
	apiKeyIndex = (apiKeyIndex + 1) % len(apollonator.APIKeys)

	payload := ApolloRequest{
		APIKey:           apollonator.APIKeys[apiKeyIndex],
		FirstName:        firstName,
		LastName:         lastName,
		OrganizationName: apollonator.Organization,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return ApolloResponse{}, err
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", apolloURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return ApolloResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ApolloResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		log.Error().Msg("API limit reached. Please try again later.")
		return ApolloResponse{}, errors.New("API limit reached")
	}

	if resp.StatusCode == 422 {
		log.Error().Msgf("API key %s is out of gas! Removing it.", apollonator.APIKeys[apiKeyIndex])

		// Remove the malfunctioning API key from the list
		apollonator.APIKeys = append(apollonator.APIKeys[:apiKeyIndex], apollonator.APIKeys[apiKeyIndex+1:]...)

		// Adjust the delay based on the new number of API keys and notify the user
		if len(apollonator.APIKeys) > 0 {
			delay = (18 * time.Second) / time.Duration(len(apollonator.APIKeys))
			log.Warn().Msgf("API key removed. Adjusting delay to %s based on the number of available API keys.", delay)
		}

		// Try the request with the next valid API key without modifying the apiKeyIndex since the list has shifted
		return ApolloRequester(apollonator, firstName, lastName, delay, apiKeyIndex)
	}

	if resp.StatusCode != http.StatusOK {
		return ApolloResponse{}, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ApolloResponse{}, err
	}

	var result ApolloResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return ApolloResponse{}, err
	}

	time.Sleep(delay)

	return result, nil
}

// Update SaveToExcel to accept a file name
func SaveToExcel(personData []PersonData, filename string) error {
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("Employee Info")

	// Adding headers
	headerRow := sheet.AddRow()
	headerRow.AddCell().Value = "First Name"
	headerRow.AddCell().Value = "Last Name"
	headerRow.AddCell().Value = "Organization"
	headerRow.AddCell().Value = "Email"
	headerRow.AddCell().Value = "Domain"
	headerRow.AddCell().Value = "Title"

	for _, data := range personData {
		row := sheet.AddRow()
		row.AddCell().Value = data.FirstName
		row.AddCell().Value = data.LastName
		row.AddCell().Value = data.Organization
		row.AddCell().Value = data.Email
		row.AddCell().Value = data.Domain
		row.AddCell().Value = data.Title
	}

	err := file.Save(filename)
	if err != nil {
		return err
	}

	return nil
}

func GetNamesFromFile(namesFile string) ([]string, error) {
	namesData, err := ioutil.ReadFile(namesFile)
	if err != nil {
		return nil, err
	}

	names := strings.Split(string(namesData), "\n")
	return names, nil
}

func main() {
	fmt.Println(`                  
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
 #@@@@@@@@@@@@@%+.                :@@@@% (pollonator) (v1.0.1)

	author: github.com/loosehose
	`)

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	parser := argparse.NewParser("apollonator", "")
	configFile := parser.String("c", "config", &argparse.Options{Required: true, Help: "Import the config.yml file with updated information."})
	namesFile := parser.String("n", "names", &argparse.Options{Required: true, Help: "Input a list of names."})
	sleep := parser.Int("s", "sleep", &argparse.Options{Default: defaultSleepDelay, Help: "Specify sleep delay in seconds. Recommended: 18 seconds per API key."})
	excelPath := parser.String("e", "excel", &argparse.Options{Help: "Specify the path and name of the excel file. If only -e is provided without a value, the default naming convention will be used."})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	apollonator, err := ParseYaml(*configFile)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse YAML config")
		os.Exit(1)
	}

	names, err := GetNamesFromFile(*namesFile)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read names file")
		os.Exit(1)
	}

	// Determine the Excel file name and path
	var excelFile string
	if *excelPath != "" {
		excelFile = *excelPath
	} else {
		re := regexp.MustCompile("[^a-zA-Z0-9_]+")
		excelFile = fmt.Sprintf("apollonator_%s.xlsx", strings.ToLower(re.ReplaceAllString(apollonator.Organization, "")))
	}

	log.Info().Msgf("Estimated completion time: %d minutes", len(names)**sleep/60)

	var personData []PersonData

	for _, line := range names {
		nameParts := strings.Fields(line)
		if len(nameParts) < 2 {
			log.Warn().Msgf("Name \"%s\" cannot be used, it must have a first name and a last name.", line)
			continue
		}
		firstName := strings.Join(nameParts[:len(nameParts)-1], " ")
		lastName := nameParts[len(nameParts)-1]

		if strings.ToLower(firstName) == "linkedin" && strings.ToLower(lastName) == "member" {
			continue
		}

		response, err := ApolloRequester(&apollonator, firstName, lastName, time.Duration(*sleep)*time.Second, apiKeyIndex)
		if err != nil {
			if err.Error() == "API limit reached" && excelFile != "" {
				err := SaveToExcel(personData, excelFile)
				if err != nil {
					log.Error().Err(err).Msg("Failed to save Excel file")
				}
				os.Exit(1)
			}
			log.Error().Err(err).Msg("Failed to get a response from Apollo")
			continue
		}

		if apollonator.Email && response.Person.Email == "" {
			continue
		}

		domain := ""
		if response.Person.Email != "" {
			domainParts := strings.Split(response.Person.Email, "@")
			domain = domainParts[1]
		}

		personData = append(personData, PersonData{
			FirstName:    firstName,
			LastName:     lastName,
			Organization: apollonator.Organization,
			Email:        response.Person.Email,
			Domain:       domain,
			Title:        response.Person.Title,
		})

		log.Info().Msgf("%s %s %s", firstName, lastName, response.Person.Email)
	}

	if excelFile != "" {
		err := SaveToExcel(personData, excelFile)
		if err != nil {
			log.Error().Err(err).Msg("Failed to save Excel file")
			os.Exit(1)
		}
	}
}
