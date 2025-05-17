package demo

import (
	"flag"
	"os"

	"github.com/krelinga/video-in-be/ffprobe"
	"github.com/krelinga/video-in-be/nfo"
	"github.com/krelinga/video-in-be/tmdb"
)

var videoFlag = flag.String("video", "", "path to the video file")

func movieNfo() error {
	const id = 170 // 28 days later.
	movieDetails, err := tmdb.GetMovieDetails(id)
	if err != nil {
		return err
	}
	fileInfo, err := ffprobe.New(*videoFlag)
	if err != nil {
		return err
	}
	movie, err := nfo.NewMovie(movieDetails, fileInfo)
	if err != nil {
		return err
	}
	err = movie.Write(os.Stdout)
	if err != nil {
		return err
	}

	return nil
}
