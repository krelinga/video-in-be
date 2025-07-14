package thumbs

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
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
	ErrFFProbe                 = errors.New("ffprobe error")
	ErrStat                    = errors.New("stat error")
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
		if err := ffmpeg(d, vf.Name(), dur); err != nil {
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

func ProjectThumbsDir(project string) string {
	return filepath.Join(env.ThumbsDir(), project)
}

func thumbsDir(d *disc) string {
	return filepath.Join(ProjectThumbsDir(d.Project), d.Disc)
}

func ffmpeg(d *disc, vf string, offset time.Duration) error {
	thumbPath := filepath.Join(env.ThumbsDir(), d.Project, d.Disc, fmt.Sprintf("%s.jpg", vf))
	videoPath := filepath.Join(discPath(d), vf)
	offsetStr := fmt.Sprintf("%02d:%02d:%02d", int(offset.Hours()), int(offset.Minutes())%60, int(offset.Seconds())%60)
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-ss", offsetStr, "-frames:v", "1", "-q:v", "2", "-y", thumbPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %w: %s", ErrFFMpeg, err, stderr.String())
	}

	probe, err := ffprobe.New(videoPath)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFFProbe, err)
	}

	stat, err := os.Stat(videoPath)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrStat, err)
	}

	found := state.ProjectModify(d.Project, func(p *state.Project) {
		disc := p.FindDiscByName(d.Disc)
		if disc == nil {
			err = errors.New("unknown disc")
			return
		}
		if disc.FindFileByName(vf) != nil {
			panic("file state already exists")
		}
		disc.Files = append(disc.Files, &state.File{
			Name:          vf,
			Thumbnail:     fmt.Sprintf("%s.jpg", vf),
			HumanByteSize: humanize.IBytes(uint64(stat.Size())),
			HumanDuration: func() string {
				if d, ok := probe.GetDuration(); ok {
					return d.String()
				}
				return ""
			}(),
			NumChapters: func() int32 {
				if ch, ok := probe.GetNumChapters(); ok {
					return ch
				}
				return 0
			}(),
		})
	})
	if !found {
		err = fmt.Errorf("%w %s", state.ErrUnknownProject, d.Project)
	}

	return err
}
