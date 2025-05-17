package ffprobe

type Raw struct {
	Streams []*RawStream `json:"streams"`
	Format  *RawFormat   `json:"format"`
}

type RawStream struct {
	CodecName     string          `json:"codec_name"`
	CodecLongName string          `json:"codec_long_name"`
	CodecType     string          `json:"codec_type"`
	Width         int64           `json:"width"`
	Height        int64           `json:"height"`
	Channels      int64           `json:"channels"`
	ChannelLayout string          `json:"channel_layout"`
	Tags          []*RawTags      `json:"tags"`
	Disposition   *RawDisposition `json:"disposition"`
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
