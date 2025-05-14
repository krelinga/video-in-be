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

func convertGenreIDs(in []int32) []string {
	// Convert genre IDs to a slice of int
	out := make([]string, 0, len(in))
	for _, id := range in {
		if name, ok := getGenre(int(id)); ok {
			out = append(out, name)
		} else {
			log.Printf("Unknown genre ID %d, skipping", id)
		}
	}
	return out
}

func convertReleaseDate(in string) (time.Time, error) {
	if in == "" {
		return time.Time{}, nil
	}
	releaseDate, err := time.Parse(time.DateOnly, in)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse release date %q: %v", in, err)
	}
	return releaseDate, nil
}

func SearchMovies(query string) ([]*MovieSearchResult, error) {
	// Search for movies
	result, err := client.SearchMovie(query, nil)
	if err != nil {
		return nil, err
	}

	out := make([]*MovieSearchResult, 0, len(result.Results))
	for _, r := range result.Results {
		releaseDate, err := convertReleaseDate(r.ReleaseDate)
		if err != nil {
			return nil, err
		}
		out = append(out, &MovieSearchResult{
			ID:            int(r.ID),
			OriginalTitle: r.OriginalTitle,
			PosterUrl:     getPosterUrl(r.PosterPath),
			Title:         r.Title,
			RealaseDate:   releaseDate,
			Overview:      r.Overview,
			Genres:        convertGenreIDs(r.GenreIDs),
		})
	}

	return out, nil
}

type MovieDetails struct {
	MovieSearchResult

	// TODO: Add more fields
}

func GetMovieDetails(id int) (*MovieDetails, error) {
	result, err := client.GetMovieInfo(id, nil)
	if err != nil {
		return nil, err
	}

	releaseDate, err := convertReleaseDate(result.ReleaseDate)
	if err != nil {
		return nil, err
	}
	genres := make([]string, 0, len(result.Genres))
	for _, g := range result.Genres {
		if g.Name == "" {
			log.Printf("Unknown genre ID %d, skipping", g.ID)
			continue
		}
		genres = append(genres, g.Name)
	}
	out := &MovieDetails{
		MovieSearchResult: MovieSearchResult{
			ID:            int(result.ID),
			OriginalTitle: result.OriginalTitle,
			PosterUrl:     getPosterUrl(result.PosterPath),
			Title:         result.Title,
			RealaseDate:   releaseDate,
			Overview:      result.Overview,
			Genres:        genres,
		},
	}
	return out, nil
}
