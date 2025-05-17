package demo

import (
	"encoding/xml"
	"os"

	"github.com/krelinga/video-in-be/ffprobe"
	"github.com/krelinga/video-in-be/nfo"
	"github.com/krelinga/video-in-be/tmdb"
)

func movieNfo() error {
	const id = 170
	movieDetails, err := tmdb.GetMovieDetails(id)
	if err != nil {
		return err
	}
	fileInfo, err := ffprobe.New("../testdata/testdata_sample_640x360.mkv")
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
