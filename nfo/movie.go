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
	Thumbs        []*Thumb
	MPAA          string `xml:"mpaa,omitempty"`
	Certification string `xml:"certification,omitempty"`
	ID            string `xml:"id,omitempty"`
	TmdbId        int    `xml:"tmdbid"`
	UniqueIds     []*UniqueId
	Genres        []string `xml:"genre"`
	Directors     []*Director
	Tags          []string `xml:"tag"`
	Actors        []*Actor
	Producers     []*Producer
	Languages     string `xml:"languages,omitempty"`
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
	Aspect       string   `xml:"aspect,omitempty"`
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

type Actor struct {
	XMLName struct{} `xml:"actor"`
	Name    string   `xml:"name"`
	Role    string   `xml:"role"`
	Thumb   string   `xml:"thumb"`
	Profile string   `xml:"profile"`
	TmdbId  int      `xml:"tmdbid"`
}

type Director struct {
	XMLName struct{} `xml:"director"`
	Name    string   `xml:",chardata"`
	TmdbId  int      `xml:"tmdbid,attr"`
}

type Producer struct {
	XMLName struct{} `xml:"producer"`
	Name    string   `xml:"name"`
	Profile string   `xml:"profile"`
	TmdbId  int      `xml:"tmdbid,attr"`
}

type Thumb struct {
	XMLName struct{} `xml:"thumb"`
	Aspect  string   `xml:"aspect,omitempty,attr"`
	URL     string   `xml:",chardata"`
}
