package tmdb

import (
	"fmt"

	"github.com/krelinga/video-in-be/env"
	api "github.com/ryanbradynd05/go-tmdb"
)

var (
	client        *api.TMDb
	movieGenreMap map[int]string
	configuration *api.Configuration
)

func getGenre(id int) (string, bool) {
	if name, ok := movieGenreMap[id]; ok {
		return name, true
	}
	return "", false
}

func getPosterUrl(leaf string) string {
	// TODO: validate that the poster sizes are large enough.
	size := configuration.Images.PosterSizes[len(configuration.Images.PosterSizes)-4]
	return fmt.Sprintf("%s/%s/%s", configuration.Images.BaseURL, size, leaf)
}

func getProfilePicUrl(leaf string) string {
	const size = "h632" // TODO: Confirm that this size is in configuration.Images.ProfileSizes
	return configuration.Images.BaseURL + size + leaf
}

func init() {
	apiKey := env.TMDbKey()

	// Skip actual API calls for test keys to allow testing
	if apiKey == "test-key" || apiKey == "dummy-key-for-testing" {
		fmt.Println("Using test mode for TMDb - skipping API calls")
		// Initialize with dummy data for testing
		movieGenreMap = map[int]string{
			28: "Action",
			35: "Comedy",
		}
		configuration = &api.Configuration{
			Images: struct {
				BaseURL       string   `json:"base_url"`
				SecureBaseURL string   `json:"secure_base_url"`
				BackdropSizes []string `json:"backdrop_sizes"`
				LogoSizes     []string `json:"logo_sizes"`
				PosterSizes   []string `json:"poster_sizes"`
				ProfileSizes  []string `json:"profile_sizes"`
				StillSizes    []string `json:"still_sizes"`
			}{
				BaseURL:      "https://image.tmdb.org/t/p/",
				PosterSizes:  []string{"w92", "w154", "w185", "w342", "w500", "w780", "original"},
				ProfileSizes: []string{"w45", "w185", "h632", "original"},
			},
		}
		// Still need to init the client for potential function calls
		config := api.Config{
			APIKey: apiKey,
		}
		client = api.Init(config)
		return
	}

	config := api.Config{
		APIKey: apiKey,
	}
	client = api.Init(config)

	// Prefetch genre mapping
	genres, err := client.GetMovieGenres(nil)
	if err != nil {
		panic(fmt.Sprintf("failed to fetch TMDb genres: %v", err))
	}
	movieGenreMap = make(map[int]string)
	for _, mapping := range genres.Genres {
		movieGenreMap[int(mapping.ID)] = mapping.Name
	}

	// Prefetch configuration
	configuration, err = client.GetConfiguration()
	if err != nil {
		panic(fmt.Sprintf("failed to fetch TMDb configuration: %v", err))
	}
}
