package demo

import (
	"context"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/krelinga/video-in-be/fanart"
	"github.com/krelinga/video-in-be/ffprobe"
	"github.com/krelinga/video-in-be/nfo"
	"github.com/krelinga/video-in-be/tmdb"
)

var directoryFlag = flag.String("dir", "", "path to the directory to refresh")
var diffNfoFlag = flag.Bool("diff_nfo", false, "print the diff between the old and new nfo files")

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
		if err := xml.Unmarshal(data, nfo); err != nil {
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
	art, err := fanart.GetArtURLs(context.Background(), oldNfo.TmdbId)
	if err != nil {
		return err
	}
	newNfo, err := nfo.NewMovie(tmdbMovie, probe, art)
	if err != nil {
		return err
	}
	var tempNfoPath string
	err = func() error {
		f, err := func() (*os.File, error) {
			if *diffNfoFlag {
				return os.CreateTemp("", "diff_nfo_*")
			} else {
				return os.Create(nfoPath)
			}
		}()
		if err != nil {
			return err
		}
		if *diffNfoFlag {
			tempNfoPath = f.Name()
		}
		defer f.Close()
		if err := newNfo.Write(f); err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		return err
	}

	if *diffNfoFlag {
		diffCmd := []string{"diff", "-u", nfoPath, tempNfoPath}
		cmd := exec.Command(diffCmd[0], diffCmd[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil && err.Error() != "exit status 1" {
			return err
		}
		if err := os.Remove(tempNfoPath); err != nil {
			return fmt.Errorf("failed to remove temp nfo file: %w", err)
		}
	}

	return nil
}
