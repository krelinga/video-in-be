package nfo

import (
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/krelinga/video-in-be/ffprobe"
	"github.com/krelinga/video-in-be/tmdb"
)

var codecMapping = func() map[string]string {
	return map[string]string{
		"mpeg2video": "MPEG-2",
		"ac3":        "AC3",
		"h264":       "h264",
	}
}()

func translateCodec(codec string) (string, error) {
	if translated, ok := codecMapping[codec]; ok {
		return translated, nil
	}
	return "", fmt.Errorf("unknown codec: %q", codec)
}

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
		UniqueIds: func() (out []*UniqueId) {
			out = [](*UniqueId){
				{
					Id:      fmt.Sprintf("%d", movieDetails.ID),
					Default: true,
					Type:    "tmdb",
				},
			}

			if movieDetails.ImdbID != "" {
				out = append(out, &UniqueId{
					Id:   movieDetails.ImdbID,
					Type: "imdb",
				})
			}

			return
		}(),
		Genres: movieDetails.Genres,
		Tags:   movieDetails.Keywords,
		FileInfo: &FileInfo{
			StreamDetails: &StreamDetails{
				Videos: slices.Collect(func(yield func(*Video) bool) {
					for stream := range probeInfo.GetVideoStreams() {
						aspect, err := stream.DisplayAspectRatio.Parse()
						if err != nil {
							log.Printf("could not parse aspect ratio %q: %v", stream.DisplayAspectRatio, err)
							aspect = ""
						}
						if !yield(&Video{
							Codec: func() string {
								c, err := translateCodec(stream.CodecName)
								if err != nil {
									setError(err)
									return ""
								}
								return c
							}(),
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
							Codec: func() string {
								c, err := translateCodec(stream.CodecName)
								if err != nil {
									setError(err)
									return ""
								}
								return c
							}(),
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
