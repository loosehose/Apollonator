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

type Config struct {
	APIKey       string `yaml:"api_key"`
	Organization string `yaml:"organization"`
	Email        bool   `yaml:"email"`
	Title        bool   `yaml:"title"`
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

// ParseYaml reads and parses a YAML configuration file
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

	// Validate Configurations
	if config.Apollonator.APIKey == "" {
		return Config{}, errors.New("missing API key in the configuration")
	}

	if config.Apollonator.Organization == "" {
		return Config{}, errors.New("missing Organization in the configuration")
	}

	return config.Apollonator, nil
}

// ApolloRequester sends a request to Apollo API and returns the response
func ApolloRequester(apollonator Config, firstName string, lastName string, delay time.Duration) (ApolloResponse, error) {
	payload := ApolloRequest{
		APIKey:           apollonator.APIKey,
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
		os.Exit(1)
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

// SaveToExcel saves the collected person data into an Excel file
func SaveToExcel(personData []PersonData, organization string) error {
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("Employee Info")

	for _, data := range personData {
		row := sheet.AddRow()
		row.AddCell().Value = data.FirstName
		row.AddCell().Value = data.LastName
		row.AddCell().Value = data.Organization
		row.AddCell().Value = data.Email
		row.AddCell().Value = data.Domain
		row.AddCell().Value = data.Title
	}

	// Using regex for edge cases such as O'Reilly
	re := regexp.MustCompile("[^a-zA-Z0-9_]+")
	safeOrganizationName := re.ReplaceAllString(organization, "")
	outputFile := fmt.Sprintf("apollonator_%s.xlsx", strings.ToLower(safeOrganizationName))

	err := file.Save(outputFile)
	if err != nil {
		return err
	}

	return nil
}

// GetNamesFromFile reads the input names file and returns a slice of names
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
	excel := parser.Flag("e", "excel", &argparse.Options{Help: "Save the results to an excel file."})
	namesFile := parser.String("n", "names", &argparse.Options{Required: true, Help: "Input a list of names."})
	sleep := parser.Int("s", "sleep", &argparse.Options{Default: defaultSleepDelay, Help: "Specify sleep delay in seconds."})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	if *configFile == "" {
		log.Error().Msg("Specify an input file [config.yml] with -c")
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

		response, err := ApolloRequester(apollonator, firstName, lastName, time.Duration(*sleep)*time.Second)
		if err != nil {
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
		// fmt.Printf("%s %s %s", firstName, lastName, response.Person.Email)
	}

	if *excel {
		err := SaveToExcel(personData, apollonator.Organization)
		if err != nil {
			log.Error().Err(err).Msg("Failed to save Excel file")
			os.Exit(1)
		}
	}
}
