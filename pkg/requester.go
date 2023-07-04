package apollonator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	apolloURL         = "https://api.apollo.io/v1/people/match"
	defaultSleepDelay = 18
)

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

// ApolloRequester sends a request to the Apollo API and returns the response
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

	if resp.StatusCode == http.StatusTooManyRequests {
		return ApolloResponse{}, fmt.Errorf("API limit reached. Please try again later")
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
