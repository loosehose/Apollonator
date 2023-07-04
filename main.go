package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/loosehose/apollonator"
)

const defaultSleepDelay = 18

func main() {
	fmt.Println(`
	Holy OSINT, Apollonator! 
	author: github.com/loosehose
	`)

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	configFile := flag.String("c", "", "Import the config.yml file with updated information.")
	excelFlag := flag.Bool("e", false, "Save the results to an excel file.")
	namesFile := flag.String("n", "", "Input a list of names.")
	sleep := flag.Int("s", defaultSleepDelay, "Specify sleep delay in seconds.")
	flag.Parse()

	if *configFile == "" {
		log.Error().Msg("Specify an input file [config.yml] with -c")
		os.Exit(1)
	}

	apollonator, err := config.ParseYaml(*configFile)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse YAML config")
		os.Exit(1)
	}

	names, err := file.GetNamesFromFile(*namesFile)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read names file")
		os.Exit(1)
	}

	log.Info().Msgf("Estimated completion time: %d minutes", len(names)**sleep/60)

	var personData []excel.PersonData

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

		response, err := requester.ApolloRequester(apollonator.APIKey, firstName, lastName, apollonator.Organization, time.Duration(*sleep)*time.Second)
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

		personData = append(personData, excel.PersonData{
			FirstName:    firstName,
			LastName:     lastName,
			Organization: apollonator.Organization,
			Email:        response.Person.Email,
			Domain:       domain,
			Title:        response.Person.Title,
		})

		log.Info().Msgf("%s %s %s", firstName, lastName, response.Person.Email)
	}

	if *excelFlag {
		err := excel.SaveToExcel(personData, apollonator.Organization)
		if err != nil {
			log.Error().Err(err).Msg("Failed to save Excel file")
			os.Exit(1)
		}
	}
}
