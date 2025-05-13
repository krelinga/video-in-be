package ffprobe

import (
	"fmt"
	"time"
)

type FFProbe struct {
	Raw *Raw
}

func (f *FFProbe) GetWidth() (int64, bool) {
	if f.Raw == nil || f.Raw.Streams == nil || len(f.Raw.Streams) == 0 {
		return 0, false
	}
	return f.Raw.Streams[0].Width, true
}

func (f *FFProbe) GetHeight() (int64, bool) {
	if f.Raw == nil || f.Raw.Streams == nil || len(f.Raw.Streams) == 0 {
		return 0, false
	}
	return f.Raw.Streams[0].Height, true
}

func (f *FFProbe) GetDuration() (time.Duration, bool) {
	if f.Raw == nil || f.Raw.Format == nil {
		return 0, false
	}
	d, err := time.ParseDuration(fmt.Sprintf("%fs", f.Raw.Format.Duration))
	if err != nil {
		panic(fmt.Sprintf("error parsing duration %f : %v", f.Raw.Format.Duration, err))
	}
	return d, true
}