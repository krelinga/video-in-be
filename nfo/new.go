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
		"dts":        "DTS",
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
		Directors: func() []*Director {
			if len(movieDetails.Crew) == 0 {
				return nil
			}
			out := make([]*Director, 0, 1)
			for _, crew := range movieDetails.Crew {
				if crew.Job == "Director" {
					out = append(out, &Director{
						Name:   crew.Name,
						TmdbId: crew.ID,
					})
				}
			}
			return out
		}(),
		Tags: movieDetails.Keywords,
		Actors: func() []*Actor {
			out := make([]*Actor, 0, len(movieDetails.Actors))
			for _, actor := range movieDetails.Actors {
				out = append(out, &Actor{
					Name:    actor.Name,
					Role:    actor.Character,
					Thumb:   actor.ProfilePicUrl,
					Profile: fmt.Sprintf("https://www.themoviedb.org/person/%d", actor.ID),
					TmdbId:  actor.ID,
				})
			}
			return out
		}(),
		Producers: func() []*Producer {
			if len(movieDetails.Crew) == 0 {
				return nil
			}
			out := make([]*Producer, 0, 5)
			for _, crew := range movieDetails.Crew {
				// TODO: this seems to catch other folks like casting ... unclear if this is correct.
				// At least it matches the current behavior of the old system.
				if crew.Department == "Production" {
					out = append(out, &Producer{
						Name:    crew.Name,
						Profile: fmt.Sprintf("https://www.themoviedb.org/person/%d", crew.ID),
						TmdbId:  crew.ID,
					})
				}
			}
			return out
		}(),
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
