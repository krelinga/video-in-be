package ffprobe

import (
	"fmt"
	"iter"
	"time"
)

type FFProbe struct {
	Raw *Raw
}

func (f *FFProbe) GetWidth() (int, bool) {
	if f.Raw == nil || f.Raw.Streams == nil || len(f.Raw.Streams) == 0 {
		return 0, false
	}
	return f.Raw.Streams[0].Width, true
}

func (f *FFProbe) GetHeight() (int, bool) {
	if f.Raw == nil || f.Raw.Streams == nil || len(f.Raw.Streams) == 0 {
		return 0, false
	}
	return f.Raw.Streams[0].Height, true
}

func (f *FFProbe) GetDuration() (time.Duration, bool) {
	if f.Raw == nil || f.Raw.Format == nil {
		return 0, false
	}
	d, err := time.ParseDuration(fmt.Sprintf("%ss", f.Raw.Format.Duration))
	if err != nil {
		panic(fmt.Sprintf("error parsing duration %s : %v", f.Raw.Format.Duration, err))
	}
	return d, true
}

func (f *FFProbe) getStreams(streamType string) iter.Seq[*RawStream] {
	return func(yield func(*RawStream) bool) {
		if f.Raw == nil || f.Raw.Streams == nil {
			return
		}
		for _, stream := range f.Raw.Streams {
			if stream.CodecType == streamType {
				if !yield(stream) {
					return
				}
			}
		}
	}
}

func (f *FFProbe) GetVideoStreams() iter.Seq[*RawStream] {
	return f.getStreams("video")
}

func (f *FFProbe) GetAudioStreams() iter.Seq[*RawStream] {
	return f.getStreams("audio")
}

func (f *FFProbe) GetSubtitleStreams() iter.Seq[*RawStream] {
	return f.getStreams("dvd_subtitle")
}
