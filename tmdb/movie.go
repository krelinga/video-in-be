package tmdb

import (
	"fmt"
	"log"
	"time"
)

type MovieSearchResult struct {
	ID            int
	OriginalTitle string
	PosterUrl     string
	Title         string
	RealaseDate   time.Time
	Overview      string
	Genres        []string
}

func SearchMovies(query string) ([]*MovieSearchResult, error) {
	// Search for movies
	result, err := client.SearchMovie(query, nil)
	if err != nil {
		return nil, err
	}

	out := make([]*MovieSearchResult, 0, len(result.Results))
	for _, r := range result.Results {
		// Get the genre names
		genreNames := make([]string, 0, len(r.GenreIDs))
		for _, id := range r.GenreIDs {
			if name, ok := getGenre(int(id)); ok {
				genreNames = append(genreNames, name)
			} else {
				log.Printf("Unknown genre ID %d, skipping", id)
			}
		}

		var releaseDate time.Time
		if r.ReleaseDate != "" {
			releaseDate, err = time.Parse(time.DateOnly, r.ReleaseDate)
			if err != nil {
				return nil, fmt.Errorf("failed to parse release date %q: %v", r.ReleaseDate, err)
			}
		}

		out = append(out, &MovieSearchResult{
			ID:            int(r.ID),
			OriginalTitle: r.OriginalTitle,
			PosterUrl:     getPosterUrl(r.PosterPath),
			Title:         r.Title,
			RealaseDate:   releaseDate,
			Overview:      r.Overview,
			Genres:        genreNames,
		})
	}

	return out, nil
}
