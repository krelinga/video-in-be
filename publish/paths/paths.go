package paths

import (
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/krelinga/video-in-be/env"
	"github.com/krelinga/video-in-be/tmdb"
)

type Paths struct {
	tmdbMovie *tmdb.MovieDetails
}

func New(tmdbMovie *tmdb.MovieDetails) Paths {
	return Paths{
		tmdbMovie: tmdbMovie,
	}
}

func (p Paths) Main(extension string) string {
	return filepath.Join(p.LibraryDir(), fmt.Sprintf("%s%s", p.name(), extension))
}

func (p Paths) Extra() string {
	randomBit := uuid.NewString()
	return filepath.Join(p.ExtrasDir(), fmt.Sprintf("%s.mkv", randomBit))
}

func (p Paths) LibraryDir() string {
	return filepath.Join(env.LibraryDir(), p.name())
}

func (p Paths) ExtrasDir() string {
	return filepath.Join(p.LibraryDir(), "extras")
}

func (p Paths) name() string {
	return fmt.Sprintf("%s (%d)", p.tmdbMovie.Title, p.tmdbMovie.ID)
}
