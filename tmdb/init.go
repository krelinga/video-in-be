package tmdb

import (
	"context"
	"fmt"
	"slices"

	api "github.com/krelinga/go-tmdb"
	"github.com/krelinga/video-in-be/env"
)

var (
	client        api.Client
	movieGenreMap map[int]string
	configuration api.ConfigDetails
)

func getGenre(id int) (string, bool) {
	if name, ok := movieGenreMap[id]; ok {
		return name, true
	}
	return "", false
}

func getPosterUrl(leaf string) string {
	// TODO: handle errors instead of returning empty string.
	const offset = -4
	if images, err := configuration.Images(); err != nil {
		return ""
	} else if baseUrl, err := images.BaseURL(); err != nil {
		return ""
	} else if posterSizes, err := images.PosterSizes(); err != nil {
		return ""
	} else if len(posterSizes) < 4 {
		return ""
	} else {
		return fmt.Sprintf("%s/%s/%s", baseUrl, posterSizes[len(posterSizes)+offset], leaf)
	}
}

func getPosterUrlOrig(leaf string) string {
	// TODO: handle errors instead of returning empty string.
	if images, err := configuration.Images(); err != nil {
		return ""
	} else if baseUrl, err := images.BaseURL(); err != nil {
		return ""
	} else {
		return fmt.Sprintf("%s%s%s", baseUrl, "original", leaf)
	}
}

func getProfilePicUrl(leaf string) string {
	// TODO: handle error instead of returning empty string.
	const size = "h632"
	if images, err := configuration.Images(); err != nil {
		return ""
	} else if profileSizes, err := images.ProfileSizes(); err != nil {
		return ""
	} else if slices.Index(profileSizes, "h632") == -1 {
		return ""
	} else if baseUrl, err := images.BaseURL(); err != nil {
		return ""
	} else if profileSizes, err := images.ProfileSizes(); err != nil {
		return ""
	} else if slices.Index(profileSizes, "h632") == -1 {
		return ""
	} else {
		return fmt.Sprintf("%s/%s/%s", baseUrl, size, leaf)
	}
}

func init() {
	apiKey := env.TMDbKey()
	client = api.ClientOptions{
		APIKey: apiKey,
	}.NewClient()

	// Skip actual API calls for test keys to allow testing
	if apiKey == "test-key" || apiKey == "dummy-key-for-testing" {
		fmt.Println("Using test mode for TMDb - skipping API calls")
		// Initialize with dummy data for testing
		movieGenreMap = map[int]string{
			28: "Action",
			35: "Comedy",
		}
		configuration = api.ConfigDetails{
			"images": api.ConfigImages{
				"base_url":      "https://image.tmdb.org/t/p/",
				"poster_sizes":  []string{"w92", "w154", "w185", "w342", "w500", "w780", "original"},
				"profile_sizes": []string{"w45", "w185", "h632", "original"},
			},
		}
		return
	}

	// Prefetch genre mapping
	movieGenreMap = make(map[int]string)
	if genres, err := api.GetMovieGenres(context.Background(), client); err != nil {
		panic(fmt.Sprintf("failed to fetch TMDb genres: %v", err))
	} else if genreList, err := genres.Genres(); err != nil {
		panic(fmt.Sprintf("failed to parse TMDb genres: %v", err))
	} else {
		for _, mapping := range genreList {
			if id, err := mapping.ID(); err != nil {
				continue
			} else if name, err := mapping.Name(); err != nil {
				continue
			} else {
				movieGenreMap[int(id)] = name
			}
		}
	}

	// Prefetch configuration
	var err error
	configuration, err = api.GetConfigDetails(context.Background(), client)
	if err != nil {
		panic(fmt.Sprintf("failed to fetch TMDb configuration: %v", err))
	}
}
