package demo

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/krelinga/video-in-be/ffprobe"
	"github.com/krelinga/video-in-be/nfo"
	"github.com/krelinga/video-in-be/tmdb"
)

var directoryFlag = flag.String("dir", "", "path to the directory to refresh")

func refresh() error {
	// Make sure the existing directory has a .nfo file.
	nfoPath, err := func() (string, error) {
		entries, err := os.ReadDir(*directoryFlag)
		if err != nil {
			return "", err
		}
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".nfo") {
				return filepath.Join(*directoryFlag, e.Name()), nil
			}
		}
		return "", errors.New("no .nfo file found")
	}()
	if err != nil {
		return err
	}

	// Read the existing .nfo file
	oldNfo, err := func() (*nfo.Movie, error) {
		// TODO: start here.
		data, err := os.ReadFile(nfoPath)
		if err != nil {
			return nil, err
		}
		nfo := &nfo.Movie{}
		if err := json.Unmarshal(data, nfo); err != nil {
			return nil, err
		}
		return nfo, nil
	}()
	if err != nil {
		return err
	}

	tmdbMovie, err := tmdb.GetMovieDetails(oldNfo.TmdbId)
	if err != nil {
		return err
	}

	// Create a new nfo file.
	mkvPath := strings.TrimSuffix(nfoPath, ".nfo") + ".mkv"
	probe, err := ffprobe.New(mkvPath)
	if err != nil {
		return err
	}
	newNfo, err := nfo.NewMovie(tmdbMovie, probe)
	if err != nil {
		return err
	}
	err = func() error {
		f, err := os.Create(nfoPath)
		if err != nil {
			return err
		}
		defer f.Close()
		if err := newNfo.Write(f); err != nil {
			return err
		}
		return nil
	}()

	return err
}
