package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/krelinga/video-in-be/tmdb"
)

func searchMovies() error {
	name := "Fight Cl"
	found, err := tmdb.SearchMovies(name)
	if err != nil {
		return fmt.Errorf("failed to search movies: %w", err)
	}

	if len(found) == 0 {
		fmt.Println("No movies found")
		return nil
	}
	for _, movie := range found {
		fmt.Printf("Found movie: %s (%s)\n", movie.Title, movie.RealaseDate.Format(time.DateOnly))
		fmt.Printf("Genres: %s\n", strings.Join(movie.Genres, ", "))
		fmt.Printf("Poster URL: %s\n", movie.PosterUrl)
		fmt.Printf("Overview: %s\n", movie.Overview)
		fmt.Println()
	}
	return nil
}
