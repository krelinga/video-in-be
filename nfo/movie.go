package nfo

type Movie struct {
	XMLName       struct{} `xml:"movie"`
	Title         string   `xml:"title"`
	OriginalTitle string   `xml:"originaltitle"`
	Year          int      `xml:"year"`
	Plot          string   `xml:"plot"`
	Outline       string   `xml:"outline"`
	Tagline       string   `xml:"tagline"`
	Runtime       int      `xml:"runtime"`
	TmdbId        int      `xml:"tmdbid"`
	UniqueIds     []*UniqueId
	Genres        []string `xml:"genre"`
	Tags          []string `xml:"tag"`
	FileInfo      *FileInfo
}

type UniqueId struct {
	XMLName struct{} `xml:"uniqueid"`
	Id      string   `xml:",chardata"`
	Default bool     `xml:"default,attr"`
	Type    string   `xml:"type,attr"`
}

type FileInfo struct {
	XMLName       struct{} `xml:"fileinfo"`
	StreamDetails *StreamDetails
}

type StreamDetails struct {
	XMLName   struct{} `xml:"streamdetails"`
	Videos    []*Video
	Audios    []*Audio
	Subtitles []*Subtitle
}

type Video struct {
	XMLName      struct{} `xml:"video"`
	Codec        string   `xml:"codec"`
	Aspect       string   `xml:"aspect"`
	Width        int      `xml:"width"`
	Height       int      `xml:"height"`
	DurationSecs int      `xml:"durationinseconds"`
}

type Audio struct {
	XMLName  struct{} `xml:"audio"`
	Codec    string   `xml:"codec"`
	Channels int      `xml:"channels"`
	Language string   `xml:"language"`
}

type Subtitle struct {
	XMLName  struct{} `xml:"subtitle"`
	Language string   `xml:"language"`
}
