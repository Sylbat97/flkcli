package flkutils

import (
	"gopkg.in/masci/flickr.v3"

	"errors"
	"flkcli/config"
	"fmt"
)

func GetFlickrClient() (*flickr.FlickrClient, error) {
	api, token, err := ValidateAndReadConfiguration()
	if err != nil {
		return nil, err
	}

	// Create a new Flickr client
	client := flickr.NewFlickrClient(api.Key, api.Secret)
	client.OAuthToken = token.OAuthToken
	client.OAuthTokenSecret = token.OAuthTokenSecret
	return client, nil
}

// ValidateConfiguration function use the config package to validate that the api key and secret are set and the oauth token and secret are set
func ValidateAndReadConfiguration() (*config.APIConfig, *config.TokenConfig, error) {
	// Get the API config
	apiConfig, err := config.GetAPIConfig()
	if err != nil {
		return nil, nil, fmt.Errorf(`failed to get API config: %w
Please use flkcli setup first
`, err)
	}

	// Check if the API key and secret are set
	if apiConfig.Key == "" || apiConfig.Secret == "" {
		return nil, nil, errors.New(`API key and secret are not set
Please use flkcli setup first
`)
	}

	// Get the token config
	tokenConfig, err := config.GetTokenConfig()
	if err != nil {
		return nil, nil, fmt.Errorf(`failed to get token config: %w
Please use flkcli login first
`, err)
	}

	// Check if the oauth token and secret are set
	if tokenConfig.OAuthToken == "" || tokenConfig.OAuthTokenSecret == "" {
		return nil, nil, errors.New(`OAuth token and secret are not set
Please use flkcli login first
`)
	}

	return apiConfig, tokenConfig, nil
}
