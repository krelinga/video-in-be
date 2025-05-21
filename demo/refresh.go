package demo

import "os"
import "strings"
import "errors"
import "github.com/krelinga/video-in-be/nfo"
import "encoding/json"
import "path/filepath"

var directoryFlag = flag.String("dir", "", "path to the directory to refresh")

func refresh() error {
	// Make sure the existing directory has a .nfo file.
	nfoPath, err := func() (string, error) {
		entries, err := os.ReadDir(*directoryFlag)
		if err != nil {
			return "", err
		}
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".nfo")
			return filepath.Join(*directoryFlag, e.Name()), nil
		}
		return "", errors.New("No .nfo file found")
	}()
	if err != nil {
		return err
	}

	// Read the existing .nfo file
	oldNfo, err := func() (*nfo.Movie, error) {
		// TODO: start here.
	}()
	return nil
}