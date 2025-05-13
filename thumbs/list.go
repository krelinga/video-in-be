package thumbs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ErrList = errors.New("could not list thumbs directory")

func List(projectName, discName string) ([]string, error) {
	thumbs, err := os.ReadDir(thumbsDir(&disc{Project: projectName, Disc: discName}))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrList, err)
	}
	out := make([]string, 0, len(thumbs))
	for _, thumb := range thumbs {
		if thumb.IsDir() || !strings.HasSuffix(thumb.Name(), ".jpg") {
			continue
		}
		out = append(out, filepath.Join(projectName, discName, thumb.Name()))
	}
	return out, nil
}