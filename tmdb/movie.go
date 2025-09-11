package tmdb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	api "github.com/krelinga/go-tmdb"
)

type MovieSearchResult struct {
	ID            int
	OriginalTitle string
	PosterUrl     string
	PosterUrlOrig string
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

func zeroDefault[T any](in T, err error) (T, error) {
	if errors.Is(err, api.ErrFieldNotFound) || errors.Is(err, api.ErrNullValue) {
		var zero T
		return zero, nil
	}
	return in, err
}

func SearchMovies(query string) ([]*MovieSearchResult, error) {
	// Search for movies
	result, err := api.SearchMovie(context.Background(), client, query)
	if err != nil {
		return nil, err
	}
	results, err := result.Results()
	if err != nil {
		return nil, err
	}
	// Map results to MovieSearchResult
	out := make([]*MovieSearchResult, 0, len(results))
	for i, r := range results {
		outResult := &MovieSearchResult{}

		if id, err := r.ID(); err != nil {
			return nil, fmt.Errorf("failed to get ID for movie at index %d: %v", i, err)
		} else {
			outResult.ID = int(id)
		}
		if originalTitle, err := zeroDefault(r.OriginalTitle()); err != nil {
			return nil, fmt.Errorf("failed to get OriginalTitle for movie at index %d: %v", i, err)
		} else {
			outResult.OriginalTitle = originalTitle
		}
		if posterPath, err := zeroDefault(r.PosterPath()); err != nil {
			return nil, fmt.Errorf("failed to get PosterPath for movie at index %d: %v", i, err)
		} else {
			outResult.PosterUrl = getPosterUrl(posterPath)
		}
		if title, err := r.Title(); err != nil {
			return nil, fmt.Errorf("failed to get Title for movie at index %d: %v", i, err)
		} else {
			outResult.Title = title
		}
		if releaseDate, err := zeroDefault(r.ReleaseDate()); err != nil {
			return nil, fmt.Errorf("failed to get ReleaseDate for movie at index %d: %v", i, err)
		} else if convertedReleaseDate, err := convertReleaseDate(releaseDate); err != nil {
			return nil, fmt.Errorf("failed to parse ReleaseDate for movie at index %d: %v", i, err)
		} else {
			outResult.RealaseDate = convertedReleaseDate
		}
		if overview, err := zeroDefault(r.Overview()); err != nil {
			return nil, fmt.Errorf("failed to get Overview for movie at index %d: %w", i, err)
		} else {
			outResult.Overview = overview
		}
		if genreIds, err := zeroDefault(r.GenreIDs()); err != nil {
			return nil, fmt.Errorf("failed to get GenreIDs for movie at index %d: %v", i, err)
		} else {
			outResult.Genres = convertGenreIDs(genreIds)
		}
		if imdbID, err := zeroDefault(r.IMDBID()); err != nil {
			return nil, fmt.Errorf("failed to get ImdbID for movie at index %d: %v", i, err)
		} else {
			outResult.ImdbID = imdbID
		}
		out = append(out, outResult)
	}

	return out, nil
}

type MovieDetails struct {
	MovieSearchResult
	Tagline    string
	Runtime    time.Duration
	Keywords   []string
	Actors     []*Actor
	Crew       []*Crew
	MPAARating string
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
	result, err := api.GetMovie(context.Background(), client, int32(id), api.WithAppendToResponse("keywords", "credits", "release_dates"))
	if err != nil {
		return nil, err
	}

	out := &MovieDetails{}
	if id, err := result.ID(); err != nil {
		return nil, fmt.Errorf("failed to get ID for movie %d: %v", id, err)
	} else {
		out.ID = int(id)
	}
	if originalTitle, err := zeroDefault(result.OriginalTitle()); err != nil {
		return nil, fmt.Errorf("failed to get original title for movie %d: %v", id, err)
	} else {
		out.OriginalTitle = originalTitle
	}
	if posterPath, err := zeroDefault(result.PosterPath()); err != nil {
		return nil, fmt.Errorf("failed to get poster path for movie %d: %v", id, err)
	} else {
		out.PosterUrl = getPosterUrl(posterPath)
		out.PosterUrlOrig = getPosterUrlOrig(posterPath)
	}
	if title, err := result.Title(); err != nil {
		return nil, fmt.Errorf("failed to get title for movie %d: %v", id, err)
	} else {
		out.Title = title
	}
	if releaseDate, err := zeroDefault(result.ReleaseDate()); err != nil {
		return nil, fmt.Errorf("failed to get release date for movie %d: %v", id, err)
	} else if parsedReleaseDate, err := convertReleaseDate(releaseDate); err != nil {
		return nil, fmt.Errorf("failed to parse release date for movie %d: %v", id, err)
	} else {
		out.RealaseDate = parsedReleaseDate
	}
	if overview, err := zeroDefault(result.Overview()); err != nil {
		return nil, fmt.Errorf("failed to get overview for movie %d: %v", id, err)
	} else {
		out.Overview = overview
	}
	if genreList, err := zeroDefault(result.Genres()); err != nil {
		return nil, fmt.Errorf("failed to get genres for movie %d: %v", id, err)
	} else {
		out.Genres = make([]string, 0, len(genreList))
		for i, g := range genreList {
			if name, err := g.Name(); err != nil {
				return nil, fmt.Errorf("failed to get name for genre ID for index %d: %v", i, err)
			} else {
				out.Genres = append(out.Genres, name)
			}
		}
	}
	if imdbId, err := zeroDefault(result.IMDBID()); err != nil {
		return nil, fmt.Errorf("failed to get IMDB ID for movie %d: %v", id, err)
	} else {
		out.ImdbID = imdbId
	}
	if tagline, err := zeroDefault(result.Tagline()); err != nil {
		return nil, fmt.Errorf("failed to get tagline for movie %d: %v", id, err)
	} else {
		out.Tagline = tagline
	}
	if runtime, err := zeroDefault(result.Runtime()); err != nil {
		return nil, fmt.Errorf("failed to get runtime for movie %d: %v", id, err)
	} else {
		out.Runtime = time.Duration(runtime) * time.Minute
	}
	if keywords, err := zeroDefault(result.Keywords()); err != nil {
		return nil, fmt.Errorf("failed to get keywords for movie %d: %v", id, err)
	} else if keywordList, err := zeroDefault(keywords.Keywords()); err != nil {
		return nil, fmt.Errorf("failed to get keyword list for movie %d: %v", id, err)
	} else {
		out.Keywords = make([]string, 0, len(keywordList))
		for i, k := range keywordList {
			if name, err := k.Name(); err != nil {
				return nil, fmt.Errorf("failed to get name for keyword ID at index %d: %v", i, err)
			} else {
				out.Keywords = append(out.Keywords, name)
			}
		}
	}
	if credits, err := zeroDefault(result.Credits()); err != nil {
		return nil, fmt.Errorf("failed to get credits for movie %d: %v", id, err)
	} else {
		if cast, err := zeroDefault(credits.Cast()); err != nil {
			return nil, fmt.Errorf("failed to get cast for movie %d: %v", id, err)
		} else {
			out.Actors = make([]*Actor, 0, len(cast))
			for i, a := range cast {
				if name, err := a.Name(); err != nil {
					return nil, fmt.Errorf("failed to get name for actor ID at index %d: %v", i, err)
				} else if character, err := zeroDefault(a.Character()); err != nil {
					return nil, fmt.Errorf("failed to get character for actor ID at index %d: %v", i, err)
				} else if profilePath, err := zeroDefault(a.ProfilePath()); err != nil {
					return nil, fmt.Errorf("failed to get profile path for actor ID at index %d: %v", i, err)
				} else if id, err := a.ID(); err != nil {
					return nil, fmt.Errorf("failed to get ID for actor ID at index %d: %v", i, err)
				} else {
					out.Actors = append(out.Actors, &Actor{
						Name:          name,
						Character:     character,
						ProfilePicUrl: getProfilePicUrl(profilePath),
						ID:            int(id),
					})
				}
			}
		}
		if crew, err := zeroDefault(credits.Crew()); err != nil {
			return nil, fmt.Errorf("failed to get crew for movie %d: %v", id, err)
		} else {
			out.Crew = make([]*Crew, 0, len(crew))
			for i, c := range crew {
				if name, err := c.Name(); err != nil {
					return nil, fmt.Errorf("failed to get name for crew ID at index %d: %v", i, err)
				} else if department, err := zeroDefault(c.Department()); err != nil {
					return nil, fmt.Errorf("failed to get department for crew ID at index %d: %v", i, err)
				} else if job, err := zeroDefault(c.Job()); err != nil {
					return nil, fmt.Errorf("failed to get job for crew ID at index %d: %v", i, err)
				} else if profilePath, err := zeroDefault(c.ProfilePath()); err != nil {
					return nil, fmt.Errorf("failed to get profile path for crew ID at index %d: %v", i, err)
				} else if id, err := c.ID(); err != nil {
					return nil, fmt.Errorf("failed to get ID for crew ID at index %d: %v", i, err)
				} else {
					out.Crew = append(out.Crew, &Crew{
						Name:          name,
						Department:    department,
						Job:           job,
						ProfilePicUrl: getProfilePicUrl(profilePath),
						ID:            int(id),
					})
				}
			}
		}
	}
	if releaseDates, err := result.ReleaseDates(); err != nil {
	} else if reults, err := releaseDates.Results(); err != nil {
	} else {
		for _, r := range reults {
			if countryCode, err := r.ISO3166_1(); err != nil {
			} else if countryCode != "US" {
			} else if releases, err := r.ReleaseDates(); err != nil {
			} else {
				for _, release := range releases {
					if cert, err := release.Certification(); err != nil {
					} else if cert == "" {
					} else {
						out.MPAARating = cert
						break
					}
				}
				break
			}
		}
	}
	return out, nil
}
