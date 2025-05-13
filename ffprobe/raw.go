package ffprobe

type Raw struct {
	Streams []*RawStream `json:"streams"`
	Format  *RawFormat   `json:"format"`
}

type RawStream struct {
	Width			int64  `json:"width"`
	Height			int64  `json:"height"`
}

type RawFormat struct {
	Duration		string `json:"duration"`
}