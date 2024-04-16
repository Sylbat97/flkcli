package list

import (
	"flkcli/cmd"
	"flkcli/flkutils"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"gopkg.in/masci/flickr.v3"
	"gopkg.in/masci/flickr.v3/photos"
	"gopkg.in/masci/flickr.v3/photosets"
)

type LoadingBar struct {
	// The total number of items to load
	Total int
	// The number of items loaded
	Loaded int
	// The last item loaded
	LastLoaded string
}

type PhotoUploadInfo struct {
	FilePath        string
	Overwrite       bool
	PreviousPhotoId string
}

var syncMode bool
var overwrite bool
var delete bool
var setName string
var setId string
var photosInSet []photosets.Photo

var UploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload photos to flickr",
	Long:  `Upload all the photos in the specified directory to flickr`,
	Run: func(command *cobra.Command, args []string) {
		sets, _ := command.Flags().GetStringArray("add-to-sets")
		syncMode, _ := command.Flags().GetBool("sync")
		overwrite, _ := command.Flags().GetBool("overwrite")
		delete, _ := command.Flags().GetBool("delete")
		if err := validateArgs(args, syncMode, overwrite, delete); err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		if len(args) > 0 {
			setName = args[0]
		}

		// Initialize the client
		client, err := cmd.GetFlickrClient()
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}

		// Get the current directoty path on the system
		dir, err := os.Getwd()
		// Check if there was an error
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		// Print dir
		fmt.Printf("Current directory: %s\n", dir)

		// Iterate over all files in dir
		files, err := os.ReadDir(dir)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}

		if syncMode {
			setId, err = flkutils.GetSetByName(client, setName, "")
			if err != nil {
				fmt.Printf("Error: %s", err)
				return
			}
			if setId != "" {
				photosInSet, _ = flkutils.GetPhotosInSet(client, setId, "")
			}
		}

		var photosToUpload []PhotoUploadInfo

		for _, file := range files {
			// Check if the file is a photo and ignore it if not
			if !isPhoto(file) {
				// Print that the file is not a photo
				fmt.Printf("\nIgnoring: %s\n", file.Name())
				continue
			}

			if syncMode {
				// Check if the photo already exists in the set
				if result, id := isPhotoInSet(file, photosInSet); result {
					if overwrite {
						// Add the file to the list of photos to upload
						photosToUpload = append(photosToUpload, PhotoUploadInfo{FilePath: file.Name(), Overwrite: true, PreviousPhotoId: id})
					}
					continue
				}
			}
			// Add the file to the list of photos to upload
			photosToUpload = append(photosToUpload, PhotoUploadInfo{FilePath: file.Name(), Overwrite: false, PreviousPhotoId: ""})

		}

		if len(photosToUpload) == 0 {
			fmt.Println("No photos to upload")
		} else {
			// Initialize the loading bar
			bar := LoadingBar{Total: len(photosToUpload), Loaded: 0}
			// Print the number of photos that will be uploaded
			printLoadingBar(bar)

			var mutex sync.Mutex
			var wg sync.WaitGroup
			maxGoroutines := 20
			guard := make(chan struct{}, maxGoroutines)

			// Iterate over all the photos to upload
			for _, pui := range photosToUpload {
				wg.Add(1)
				guard <- struct{}{}
				go func(pui PhotoUploadInfo) {
					defer wg.Done()

					client, err := flkutils.GetFlickrClient()
					if err != nil {
						fmt.Printf("Error: %s", err)
						<-guard
					}

					resp, err := uploadPhoto(client, pui)
					if err != nil {
						fmt.Printf("Error: %s", err)
						<-guard
					}

					for _, set := range sets {
						err := flkutils.AddToPhotoSet(client, resp.ID, set)
						if err != nil {
							fmt.Printf("Error: %s", err)
						}
					}

					mutex.Lock()
					if setId == "" && syncMode {
						setId, err = flkutils.CreateSet(client, setName, "", resp.ID)
						if err != nil {
							fmt.Printf("Error: %s", err)
						}
					} else if syncMode {
						err := flkutils.AddToPhotoSet(client, resp.ID, setId)
						if err != nil {
							fmt.Printf("Error: %s", err)
						}
					}
					bar.Loaded++
					bar.LastLoaded = pui.FilePath
					fmt.Printf("\nUploaded: %s\n", pui.FilePath)
					printLoadingBar(bar)
					mutex.Unlock()
					<-guard
				}(pui)
			}
			wg.Wait()
		}

		var photosToDelete []string

		if delete {
			for _, photo := range photosInSet {
				found := false
				for _, file := range files {
					extension := filepath.Ext(file.Name())
					fileName := strings.TrimSuffix(file.Name(), extension)
					if strings.EqualFold(strings.ToLower(photo.Title), strings.ToLower(fileName)) {
						found = true
						break
					}
				}
				if !found {
					photosToDelete = append(photosToDelete, photo.Id)
				}
			}
			if len(photosToDelete) > 0 {
				fmt.Printf("\nDeleting %d photos...\n", len(photosToDelete))
				// Ask for confirmation
				fmt.Print("Are you sure you want to delete these photos? (y/n): ")
				var response string
				fmt.Scanln(&response)
				if strings.ToLower(response) == "y" {
					for _, photoId := range photosToDelete {
						_, err := photos.Delete(client, photoId)
						if err != nil {
							fmt.Printf("Error: %s", err)
						}
					}
				}
			}
		}
	},
}

func isPhotoInSet(file fs.DirEntry, photosInSet []photosets.Photo) (bool, string) {
	for _, photo := range photosInSet {
		extension := filepath.Ext(file.Name())
		fileName := strings.TrimSuffix(file.Name(), extension)
		if strings.EqualFold(strings.ToLower(photo.Title), strings.ToLower(fileName)) {
			// Print that the file already exists in the set
			return true, photo.Id
		}
	}
	return false, ""
}

func init() {
	// Initialize empty string array
	UploadCmd.PersistentFlags().StringArray("add-to-sets", []string{}, "IDs of sets to upload the photos to")
	UploadCmd.PersistentFlags().BoolVarP(&syncMode, "sync", "s", false, "Sync mode")
	UploadCmd.PersistentFlags().BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite already existing photos")
	UploadCmd.PersistentFlags().BoolVarP(&delete, "delete", "d", false, "Delete distant photos from flickr that are not in the local directory")

	cmd.RootCmd.AddCommand(UploadCmd)
}

func validateArgs(args []string, syncMode bool, overwrite bool, delete bool) error {
	if (len(args) == 0 || args[0] == "") && syncMode {
		return fmt.Errorf("album name is required in sync mode")
	}
	if overwrite && delete && !syncMode {
		return fmt.Errorf("overwrite and delete flags can only be used in sync mode")
	}
	return nil
}

func uploadPhoto(client *flickr.FlickrClient, pui PhotoUploadInfo) (*flickr.UploadResponse, error) {
	params := flickr.NewUploadParams()
	// Restrict the photo to private by default
	params.IsPublic = false
	params.IsFamily = false
	params.IsFriend = false
	resp, err := flickr.UploadFile(client, pui.FilePath, params)
	if err != nil {
		return resp, err
	}
	if resp.Status != "ok" {
		fmt.Printf("Error: %s", resp.Extra)
		return resp, err
	}
	if pui.Overwrite {
		photos.Delete(client, pui.PreviousPhotoId)
	}
	return resp, nil
}

func isPhoto(file os.DirEntry) bool {
	// Check if the file has a valid photo extension
	extensions := []string{".jpg", ".jpeg", ".png", ".gif"}
	for _, ext := range extensions {
		if strings.HasSuffix(strings.ToLower(file.Name()), ext) {
			return true
		}
	}
	return false
}

func printLoadingBar(bar LoadingBar) {
	progress := float64(bar.Loaded) / float64(bar.Total)
	barLength := 30
	filledLength := int(progress * float64(barLength))

	barstr := strings.Repeat("█", filledLength) + strings.Repeat(" ", barLength-filledLength)

	// Clear console content
	fmt.Print("\033[H\033[2J")
	// Print total number of files being uploaded
	fmt.Printf("Uploading %d files...\n", bar.Total)
	fmt.Printf("\r[%s] %d/%d", barstr, bar.Loaded, bar.Total)
	// Print last successfully uploaded file if not empty
	if bar.LastLoaded != "" {
		fmt.Printf("\nLast uploaded: %s", bar.LastLoaded)
	}
}
