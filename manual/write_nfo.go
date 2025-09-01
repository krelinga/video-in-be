package manual

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/krelinga/video-in-be/ffprobe"
	"github.com/krelinga/video-in-be/nfo"
	"github.com/krelinga/video-in-be/tmdb"
)

func mkvPath(movieDir string) (string, error) {
	var mkvPath string
	files, err := os.ReadDir(movieDir)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) == ".mkv" {
			if mkvPath != "" {
				return "", errors.New("multiple mkv files found")
			}
			mkvPath = filepath.Join(movieDir, file.Name())
		}
	}
	if mkvPath == "" {
		return "", errors.New("no mkv file found")
	}
	return mkvPath, nil
}

func nfoPath(moviePath string) string {
	return moviePath[:len(moviePath)-len(filepath.Ext(moviePath))] + ".nfo"
}

func WriteNfo() {
	if err := writeNfo(); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing NFO: %v\n", err)
		os.Exit(1)
	}
}

func writeNfo() error {
	if *movieDirFlag == "" {
		return errors.New("movie_dir flag is required")
	}

	mkvPath, err := mkvPath(*movieDirFlag)
	if err != nil {
		return err
	}
	nfoPath := nfoPath(mkvPath)

	var movieId int
	if *tmdbMovieIdFlag != 0 {
		movieId = *tmdbMovieIdFlag
	} else {
		movieId, err = readTmdbIdFromFile(nfoPath)
		if err != nil {
			return fmt.Errorf("failed to read tmdbid from existing nfo, pass --tmdb_movie_id to bypass this error: %w", err)
		}
	}

	details, err := tmdb.GetMovieDetails(movieId)
	if err != nil {
		return err
	}
	probeInfo, err := ffprobe.New(mkvPath)
	if err != nil {
		return err
	}
	movieNfo, err := nfo.NewMovie(details, probeInfo)
	if err != nil {
		return err
	}
	var writer io.Writer
	if !*dryRunFlag {
		fileWriter, err := os.Create(nfoPath)
		if err != nil {
			return err
		}
		defer fileWriter.Close()
		writer = fileWriter
	} else {
		writer = os.Stdout
	}
	return movieNfo.Write(writer)
}
