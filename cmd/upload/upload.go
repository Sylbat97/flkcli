package list

import (
	"flkcli/cmd"
	"flkcli/flkutils"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"gopkg.in/masci/flickr.v3"
)

type LoadingBar struct {
	// The total number of items to load
	Total int
	// The number of items loaded
	Loaded int
	// The last item loaded
	LastLoaded string
}

// TODO
// Add a flag to specify the directory to upload photos from if not the current directory
// Add ability to specify a new set to create and upload photos to

var UploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload photos to flickr",
	Long:  `Upload all the photos in the specified directory to flickr`,
	Run: func(command *cobra.Command, args []string) {
		sets, _ := command.Flags().GetStringArray("sets")
		// Initialize the client
		_, err := cmd.GetFlickrClient()
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

		var photosToUpload []string

		for _, file := range files {
			// Check if the file is a photo and ignore it if not
			if !isPhoto(file) {
				// Print that the file is not a photo
				fmt.Printf("\nIgnoring: %s\n", file.Name())
				continue
			}
			// Add the file to the list of photos to upload
			photosToUpload = append(photosToUpload, file.Name())

		}

		// Initialize the loading bar
		bar := LoadingBar{Total: len(photosToUpload), Loaded: 0}
		// Print the number of photos that will be uploaded
		printLoadingBar(bar)
		var mutex sync.Mutex
		var wg sync.WaitGroup

		// Iterate over all the photos to upload
		for _, filePath := range photosToUpload {
			wg.Add(1)
			go func(filePath string) {
				defer wg.Done()

				client, err := flkutils.GetFlickrClient()
				if err != nil {
					fmt.Printf("Error: %s", err)
					return
				}

				resp, err := flickr.UploadFile(client, filePath, nil)
				if err != nil {
					fmt.Printf("Error: %s", err)
					return
				}

				if resp.Status != "ok" {
					fmt.Printf("Error: %s", resp.Extra)
					return
				}

				// Iterate over all the sets to upload the photo to
				for _, set := range sets {
					err := flkutils.AddToPhotoSet(client, resp.ID, set)
					if err != nil {
						fmt.Printf("Error: %s", err)
					}
				}

				mutex.Lock()
				bar.Loaded++
				bar.LastLoaded = filePath
				fmt.Printf("\nUploaded: %s\n", filePath)
				printLoadingBar(bar)
				mutex.Unlock()
			}(filePath)
		}

		wg.Wait()

	},
}

func init() {
	// Initialize empty string array
	UploadCmd.PersistentFlags().StringArray("sets", []string{}, "Sets to upload the photos to")
	cmd.RootCmd.AddCommand(UploadCmd)
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
