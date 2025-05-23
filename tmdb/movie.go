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
	ImdbID        string
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
	Tagline  string
	Runtime  time.Duration
	Keywords []string
	Actors   []*Actor
	Crew     []*Crew
}

type Actor struct {
	Name          string
	Character     string
	ProfilePicUrl string
	ID            int
}

type Crew struct {
	Name          string
	Department    string
	Job           string
	ProfilePicUrl string
	ID            int
}

func GetMovieDetails(id int) (*MovieDetails, error) {
	options := map[string]string{
		"append_to_response": "keywords,credits",
	}
	result, err := client.GetMovieInfo(id, options)
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
			ImdbID:        result.ImdbID,
		},
		Tagline: result.Tagline,
		Runtime: time.Duration(result.Runtime) * time.Minute,
		Keywords: func() []string {
			if result.Keywords == nil {
				return nil
			}
			out := make([]string, 0, len(result.Keywords.Keywords))
			for _, k := range result.Keywords.Keywords {
				if k.Name == "" {
					log.Printf("Unknown keyword ID %d, skipping", k.ID)
					continue
				}
				out = append(out, k.Name)
			}
			return out
		}(),
		Actors: func() []*Actor {
			if result.Credits == nil {
				return nil
			}
			out := make([]*Actor, 0, len(result.Credits.Cast))
			for _, a := range result.Credits.Cast {
				if a.Name == "" {
					log.Printf("Unknown actor ID %d, skipping", a.ID)
					continue
				}
				out = append(out, &Actor{
					Name:          a.Name,
					Character:     a.Character,
					ProfilePicUrl: getProfilePicUrl(a.ProfilePath),
					ID:            a.ID,
				})
			}
			return out
		}(),
		Crew: func() []*Crew {
			if result.Credits == nil {
				return nil
			}
			out := make([]*Crew, 0, len(result.Credits.Crew))
			for _, c := range result.Credits.Crew {
				if c.Name == "" {
					log.Printf("Unknown crew ID %d, skipping", c.ID)
					continue
				}
				out = append(out, &Crew{
					Name:          c.Name,
					Department:    c.Department,
					Job:           c.Job,
					ProfilePicUrl: getProfilePicUrl(c.ProfilePath),
					ID:            c.ID,
				})
			}
			return out
		}(),
	}
	return out, nil
}
