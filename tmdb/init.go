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

func init() {
	config := api.Config{
		APIKey: env.TMDbKey(),
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
