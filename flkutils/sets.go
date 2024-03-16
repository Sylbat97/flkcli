package flkutils

import (
	"fmt"

	"gopkg.in/masci/flickr.v3"
	"gopkg.in/masci/flickr.v3/photosets"
)

func ListSets(client *flickr.FlickrClient, userid string) (total int, photosetitems []photosets.Photoset, error error) {
	// Get the list of sets
	response, err := photosets.GetList(client, true, userid, 0)

	if err != nil {
		return 0, nil, fmt.Errorf("failed to get list of sets: %w", err)
	}

	photoSetItems := response.Photosets.Items

	for response.Photosets.Page < response.Photosets.Pages {
		response, err = photosets.GetList(client, true, userid, response.Photosets.Page+1)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to get list of sets: %w", err)
		}
		photoSetItems = append(photoSetItems, response.Photosets.Items...)
	}

	return response.Photosets.Total, photoSetItems, nil
}
