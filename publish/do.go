package publish

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/krelinga/video-in-be/ffprobe"
	"github.com/krelinga/video-in-be/nfo"
	"github.com/krelinga/video-in-be/publish/paths"
	"github.com/krelinga/video-in-be/state"
	"github.com/krelinga/video-in-be/thumbs"
	"github.com/krelinga/video-in-be/tmdb"
)

func checkProject(project *state.Project) error {
	if project.TmdbId == 0 {
		return errors.New("no TMDB ID")
	}
	mainTitleCount := 0
	for _, disc := range project.Discs {
		for _, file := range disc.Files {
			switch file.Category {
			case state.FileCatExtra: // Nothing to do
			case state.FileCatMainTitle:
				mainTitleCount++
			case state.FileCatTrash: // Nothing to do
			case state.FileCatNone:
				return fmt.Errorf("file category is not set %s/%s", disc.Name, file.Name)
			default:
				return fmt.Errorf("unknown file category %s/%s: %s", disc.Name, file.Name, file.Category)

			}
		}
	}
	if mainTitleCount != 1 {
		return fmt.Errorf("there should be exactly one main title, found %d", mainTitleCount)
	}
	return nil
}

func Do(project *state.Project) error {
	if err := checkProject(project); err != nil {
		return err
	}

	tmdbMovie, err := tmdb.GetMovieDetails(project.TmdbId)
	if err != nil {
		return err
	}

	p := paths.New(tmdbMovie)

	// Create directories
	for _, dir := range []string{p.LibraryDir(), p.ExtrasDir()} {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			if os.IsExist(err) {
				continue
			}
			return err
		}
	}

	// Figure out renames to be done.
	type rename struct {
		inPath  string
		outPath string
	}
	renames := []rename{}
	var mainPath string
	for _, disc := range project.Discs {
		for _, file := range disc.Files {
			inPath := filepath.Join(state.ProjectDir(project.Name), disc.Name, file.Name)
			var outPath string
			switch file.Category {
			case state.FileCatMainTitle:
				outPath = p.Main(".mkv")
			case state.FileCatExtra:
				outPath = p.Extra()
			default:
				continue
			}
			renames = append(renames, rename{
				inPath:  inPath,
				outPath: outPath,
			})
			if file.Category == state.FileCatMainTitle {
				mainPath = outPath
			}
		}
	}
	if mainPath == "" {
		panic("no main path")
	}

	// Generate NFO
	probeInfo, err := ffprobe.New(mainPath)
	if err != nil {
		return err
	}
	movieNfo, err := nfo.NewMovie(tmdbMovie, probeInfo)
	if err != nil {
		return err
	}
	err = func() error {
		f, err := os.Create(p.Main(".nfo"))
		if err != nil {
			return err
		}
		defer f.Close()
		if err := movieNfo.Write(f); err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		return err
	}

	// TODO: generate .tcprofile file.

	// Execute renames
	for _, rename := range renames {
		if err := os.Rename(rename.inPath, rename.outPath); err != nil {
			return fmt.Errorf("could not move %s to %s: %w", rename.inPath, rename.outPath, err)
		}
	}

	// Remove the project directory if nothing failed.
	if err := os.RemoveAll(state.ProjectDir(project.Name)); err != nil {
		return fmt.Errorf("could not remove project directory %s: %w", state.ProjectDir(project.Name), err)
	}

	// Remove thumbs
	thumbsDir := thumbs.ProjectThumbsDir(project.Name)
	if err := os.RemoveAll(thumbsDir); err != nil {
		return fmt.Errorf("could not remove thumbs directory %s: %w", thumbsDir, err)
	}

	return nil
}
