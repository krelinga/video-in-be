package thumbs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/krelinga/video-in-be/env"
	"github.com/krelinga/video-in-be/ffprobe"
	"github.com/krelinga/video-in-be/state"
)

func generateThumbs(in <-chan *disc) {
	for d := range in {
		if err := generateDiscThumbs(d); err != nil {
			trySetError(d.Project, d.Disc)
			log.Print(err)
			continue
		}
	}
}

var (
	ErrNoDuration              = errors.New("no duration found")
	ErrCouldNotReadDiscDir     = errors.New("could not read disc directory")
	ErrCouldNotCreateThumbsDir = errors.New("could not create thumbs directory")
	ErrFFMpeg                  = errors.New("ffmpeg error")
)

func generateDiscThumbs(d *disc) error {
	if err := setState(d.Project, d.Disc, state.ThumbStateWaiting, state.ThumbStateWorking); err != nil {
		return err
	}

	videoFiles, err := os.ReadDir(discPath(d))
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCouldNotReadDiscDir, err)
	}

	if err := os.MkdirAll(thumbsDir(d), 0755); err != nil {
		return fmt.Errorf("%w: %w", ErrCouldNotCreateThumbsDir, err)
	}

	for _, vf := range videoFiles {
		if vf.IsDir() || !strings.HasSuffix(vf.Name(), ".mkv") {
			continue
		}
		dur, err := getThumbOffset(d, vf.Name())
		if err != nil {
			return err
		}
		if err := ffmpeg(d, dur); err != nil {
			return err
		}
	}

	if err := setState(d.Project, d.Disc, state.ThumbStateWorking, state.ThumbStateDone); err != nil {
		return err
	}

	return nil
}

func discPath(d *disc) string {
	return filepath.Join(state.ProjectDir(d.Project), d.Disc)
}

func getThumbOffset(d *disc, vf string) (time.Duration, error) {
	info, err := ffprobe.New(filepath.Join(discPath(d), vf))
	if err != nil {
		return 0, err
	}
	dur, ok := info.GetDuration()
	if !ok {
		return 0, fmt.Errorf("%w: %s/%s", ErrNoDuration, d.Project, d.Disc)
	}
	return dur / 2, nil
}

func thumbsDir(d *disc) string {
	return filepath.Join(env.ThumbsDir(), d.Project, d.Disc)
}

func ffmpeg(d *disc, offset time.Duration) error {
	thumbPath := filepath.Join(state.ProjectDir(d.Project), d.Disc, "thumb.jpg")
	offsetStr := fmt.Sprintf("%02d:%02d:%02d", int(offset.Hours()), int(offset.Minutes())%60, int(offset.Seconds())%60)
	cmd := exec.Command("ffmpeg", "-i", thumbPath, "-ss", offsetStr, "-frames:v", "1", "-q:v", "2", "-y", thumbPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %w", ErrFFMpeg, err)
	}
	return nil
}
