package nfo

import (
	"errors"
	"fmt"
	"slices"

	"github.com/krelinga/video-in-be/ffprobe"
	"github.com/krelinga/video-in-be/tmdb"
)

func NewMovie(movieDetails *tmdb.MovieDetails, probeInfo *ffprobe.FFProbe) (outMovie *Movie, outError error) {
	setError := func(err error) {
		if outError == nil {
			outError = err
		}
	}
	// Create a new Movie
	outMovie = &Movie{
		Title:         movieDetails.Title,
		OriginalTitle: movieDetails.OriginalTitle,
		Year:          movieDetails.RealaseDate.Year(),
		Plot:          movieDetails.Overview,
		Outline:       movieDetails.Overview,
		Tagline:       movieDetails.Tagline,
		Runtime:       int(movieDetails.Runtime.Minutes()),
		TmdbId:        int(movieDetails.ID),
		UniqueIds: []*UniqueId{
			{
				Id:      fmt.Sprintf("%d", movieDetails.ID),
				Default: true,
				Type:    "tmdb",
			},
		},
		Genres: movieDetails.Genres,
		// TODO: Add the tags
		FileInfo: &FileInfo{
			StreamDetails: &StreamDetails{
				Videos: slices.Collect(func(yield func(*Video) bool) {
					for stream := range probeInfo.GetVideoStreams() {
						aspect, err := stream.DisplayAspectRatio.Parse()
						if err != nil {
							setError(err)
							return
						}
						if !yield(&Video{
							Codec:  stream.CodecLongName, // TODO: some translation needed
							Aspect: aspect,
							Width:  stream.Width,
							Height: stream.Height,
							DurationSecs: func() int {
								if d, ok := probeInfo.GetDuration(); ok {
									return int(d.Seconds())
								}
								setError(errors.New("could not get duration from ffprobe"))
								return 0
							}(),
						}) {
							return
						}
					}
				}),
				Audios: slices.Collect(func(yield func(*Audio) bool) {
					for stream := range probeInfo.GetAudioStreams() {
						if !yield(&Audio{
							Codec:    stream.CodecName, // TODO: some translation needed
							Channels: stream.Channels,
							Language: stream.Tags.Language,
						}) {
							return
						}
					}
				}),
				Subtitles: slices.Collect(func(yield func(*Subtitle) bool) {
					for stream := range probeInfo.GetAudioStreams() {
						if !yield(&Subtitle{
							Language: stream.Tags.Language,
						}) {
							return
						}
					}
				}),
			},
		},
	}

	if outError != nil {
		outMovie = nil
	}
	return
}
