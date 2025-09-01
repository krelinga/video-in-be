package manual

import (
	"encoding/xml"
	"errors"
	"os"
)

func readTmdbIdFromFile(nfoPath string) (int, error) {
	data, err := os.ReadFile(nfoPath)
	if err != nil {
		return 0, err
	}

	var nfo struct {
		Movie *struct {
			TmdbID *int `xml:"tmdbid"`
		} `xml:"movie"`
	}
	if err := xml.Unmarshal(data, &nfo); err != nil {
		return 0, err
	}
	if nfo.Movie == nil {
		return 0, errors.New("movie element is missing")
	}
	if nfo.Movie.TmdbID == nil {
		return 0, errors.New("tmdbid element is missing")
	}
	return *nfo.Movie.TmdbID, nil
}
