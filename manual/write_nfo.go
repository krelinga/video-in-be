package manual

import (
	"errors"
	"fmt"
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
	if *tmdbMovieIdFlag == 0 {
		return errors.New("tmdb_movie_id flag is required")
	}

	details, err := tmdb.GetMovieDetails(*tmdbMovieIdFlag)
	if err != nil {
		return err
	}
	mkvPath, err := mkvPath(*movieDirFlag)
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
	writer, err := os.Create(nfoPath(mkvPath))
	if err != nil {
		return err
	}
	defer writer.Close()
	return movieNfo.Write(writer)
}
