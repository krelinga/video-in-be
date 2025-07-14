package ffprobe

import (
	"fmt"
	"strconv"
	"strings"
)

type Raw struct {
	Streams  []*RawStream  `json:"streams"`
	Format   *RawFormat    `json:"format"`
	Chapters []*RawChapter `json:"chapters"`
}

type RawStream struct {
	CodecName          string                `json:"codec_name"`
	CodecLongName      string                `json:"codec_long_name"`
	CodecType          string                `json:"codec_type"`
	Width              int                   `json:"width"`
	Height             int                   `json:"height"`
	DisplayAspectRatio RawDisplayAspectRatio `json:"display_aspect_ratio"`
	Channels           int                   `json:"channels"`
	ChannelLayout      string                `json:"channel_layout"`
	Tags               *RawTags              `json:"tags"`
	Disposition        *RawDisposition       `json:"disposition"`
}

type RawDisplayAspectRatio string

func (r RawDisplayAspectRatio) Parse() (string, error) {
	parsed, ok := func() (string, bool) {
		parts := strings.Split(string(r), ":")
		if len(parts) != 2 {
			return "", false
		}
		lhsInt, err := strconv.Atoi(parts[0])
		if err != nil {
			return "", false
		}
		rhsInt, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", false
		}
		return strconv.FormatFloat(float64(lhsInt)/float64(rhsInt), 'f', 2, 64), true
	}()

	if !ok {
		return "", fmt.Errorf("invalid display aspect ratio %q", string(r))
	}
	return parsed, nil
}

type RawTags struct {
	Language string `json:"language"`
	Title    string `json:"title"`
}

type RawDisposition struct {
	Default int64 `json:"default"`
	Forced  int64 `json:"forced"`
}

type RawFormat struct {
	Duration string `json:"duration"`
}

type RawChapter struct {
	StartTime string          `json:"start_time"`
	EndTime   string          `json:"end_time"`
	Tags      *RawChapterTags `json:"tags"`
}

type RawChapterTags struct {
	Title string `json:"title"`
}
