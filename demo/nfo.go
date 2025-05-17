package demo

import (
	"encoding/xml"
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
	movie := nfo.NewMovie(movieDetails, fileInfo)
	xml, err := xml.MarshalIndent(movie, "", "  ")
	if err != nil {
		return err
	}
	os.Stdout.Write(xml)
	os.Stdout.Write([]byte("\n"))

	return nil
}
