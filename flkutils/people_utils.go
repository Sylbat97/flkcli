package flkutils

import (
	"fmt"

	"gopkg.in/masci/flickr.v3"
	"gopkg.in/masci/flickr.v3/people"
)

func ResolveId(client *flickr.FlickrClient, username string) (id string, error error) {
	user_response, err := people.FindByUsername(client, username)
	if err != nil {
		return "", fmt.Errorf("failed to resolve username: %w", err)
	}
	return user_response.User.Id, nil
}
