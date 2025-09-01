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
		TmdbID *int `xml:"tmdbid"`
	}
	if err := xml.Unmarshal(data, &nfo); err != nil {
		return 0, err
	}
	if nfo.TmdbID == nil {
		return 0, errors.New("tmdbid element is missing")
	}
	return *nfo.TmdbID, nil
}
