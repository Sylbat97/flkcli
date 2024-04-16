package flkutils

import (
	"fmt"
	"strings"

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

func AddToPhotoSet(client *flickr.FlickrClient, photoId, photoSetId string) error {
	// Get the list of sets
	response, err := photosets.AddPhoto(client, photoSetId, photoId)
	if err != nil {
		return fmt.Errorf("failed to add picture to set: %w", err)
	}

	if response.Status != "ok" {
		return fmt.Errorf("failed to add picture to set: %s", response.Extra)
	}

	return nil
}

func GetSetByName(client *flickr.FlickrClient, setName, userid string) (string, error) {
	_, sets, err := ListSets(client, userid)
	if err != nil {
		return "", fmt.Errorf("failed to get list of sets: %w", err)
	}
	for _, set := range sets {
		if strings.EqualFold(strings.ToLower(set.Title), strings.ToLower(setName)) {
			return set.Id, nil
		}
	}
	return "", nil
}

func GetPhotosInSet(client *flickr.FlickrClient, setId, userId string) (photosetitems []photosets.Photo, error error) {
	response, err := photosets.GetPhotos(client, true, setId, userId, 0)

	if err != nil {
		return nil, fmt.Errorf("failed to get list of sets: %w", err)
	}

	photoItems := response.Photoset.Photos

	for response.Photoset.Page < response.Photoset.Pages {
		response, err = photosets.GetPhotos(client, true, setId, userId, response.Photoset.Page+1)
		if err != nil {
			return nil, fmt.Errorf("failed to get list of sets: %w", err)
		}
		photoItems = append(photoItems, response.Photoset.Photos...)
	}
	return photoItems, nil
}

func CreateSet(client *flickr.FlickrClient, title, description, primaryPhotoId string) (string, error) {
	response, err := photosets.Create(client, title, description, primaryPhotoId)
	if err != nil {
		return "", fmt.Errorf("failed to create set: %w", err)
	}

	if response.Status != "ok" {
		return "", fmt.Errorf("failed to create set: %s", response.Extra)
	}

	return response.Set.Id, nil
}
