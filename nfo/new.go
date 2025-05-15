package nfo

import (
	"fmt"

	"github.com/krelinga/video-in-be/ffprobe"
	"github.com/krelinga/video-in-be/tmdb"
)

func NewMovie(movieDetails *tmdb.MovieDetails, ffprobe *ffprobe.FFProbe) *Movie {
	// Create a new Movie
	movie := &Movie{
		Title:         movieDetails.Title,
		OriginalTitle: movieDetails.OriginalTitle,
		Year:          movieDetails.RealaseDate.Year(),
		Plot:          movieDetails.Overview,
		Outline:       movieDetails.Overview,
		Tagline:       movieDetails.Tagline,
		Runtime:       int(movieDetails.Runtime.Minutes()),
		TmdbId:        int(movieDetails.ID),
		UniqueIds: []*UniqueId{
			{
				Id:      fmt.Sprintf("%d", movieDetails.ID),
				Default: true,
				Type:    "tmdb",
			},
		},
		Genres: movieDetails.Genres,
		// TODO: Add the tags
		FileInfo: &FileInfo{
			StreamDetails: &StreamDetails{},
		},
	}

	// TODO: Add the file info
	return movie
}
