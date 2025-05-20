package state

import (
	"path/filepath"
	"sync"

	"github.com/krelinga/video-in-be/env"
)

type ThumbState string

const (
	ThumbStateNone    ThumbState = ""
	ThumbStateWaiting ThumbState = "waiting"
	ThumbStateWorking ThumbState = "working"
	ThumbStateDone    ThumbState = "done"
	ThumbStateError   ThumbState = "error"
)

type FileCat string

const (
	FileCatNone      FileCat = ""
	FileCatMainTitle FileCat = "main_title"
	FileCatExtra     FileCat = "extra"
	FileCatTrash     FileCat = "trash"
)

type Project struct {
	Name   string  `json:"name"`
	Discs  []*Disc `json:"discs,omitempty"`
	TmdbId int     `json:"tmdb_id,omitempty"`
}

func (p *Project) FindDiscByName(name string) *Disc {
	for _, d := range p.Discs {
		if d.Name == name {
			return d
		}
	}
	return nil
}

type Disc struct {
	Name       string     `json:"name"`
	ThumbState ThumbState `json:"thumb_state,omitempty"`
	Files      []*File    `json:"files,omitempty"`
}

func (d *Disc) FindFileByName(name string) *File {
	for _, f := range d.Files {
		if f.Name == name {
			return f
		}
	}
	return nil
}

func (d *Disc) FindFileByThumbnail(thumbnail string) *File {
	for _, f := range d.Files {
		if f.Thumbnail == thumbnail {
			return f
		}
	}
	return nil
}

type File struct {
	Name          string  `json:"name"`
	Category      FileCat `json:"category,omitempty"`
	Thumbnail     string  `json:"thumbnail,omitempty"`
	HumanByteSize string  `json:"human_byte_size,omitempty"`
	HumanDuration string  `json:"human_duration,omitempty"`
	NumChapters   int32   `json:"num_chapters,omitempty"`
}

var projectMutex = &sync.RWMutex{}

func projectPath() string {
	return env.StateDir() + "/projects.json"
}

func ProjectsRead(fn func([]*Project)) {
	projectMutex.RLock()
	defer projectMutex.RUnlock()
	fn(read[Project](projectPath()))
}

func ProjectsModify(fn func([]*Project) []*Project) {
	projectMutex.Lock()
	defer projectMutex.Unlock()
	projects := read[Project](projectPath())
	projects = fn(projects)
	write(projectPath(), projects)
}

func ProjectDir(project string) string {
	return filepath.Join(env.ProjectDir(), project)
}

func ProjectRead(name string, fn func(*Project)) bool {
	var found bool
	ProjectsRead(func(all []*Project) {
		for _, x := range all {
			if x.Name == name {
				found = true
				if fn != nil {
					fn(x)
				}
				return
			}
		}
	})
	return found
}

func ProjectModify(name string, fn func(*Project)) bool {
	var found bool
	ProjectsModify(func(all []*Project) []*Project {
		for _, x := range all {
			if x.Name == name {
				found = true
				if fn != nil {
					fn(x)
				}
				break
			}
		}
		return all
	})
	return found
}

func ProjectReadAndRemove(name string, fn func(*Project) error) bool {
	var found bool
	ProjectsModify(func(in []*Project) []*Project {
		out := make([]*Project, 0, len(in))
		for _, x := range in {
			if x.Name == name {
				found = true
				err := fn(x)
				if err != nil {
					// Only remove this project if the function returns nil
					return in
				}
			} else {
				out = append(out, x)
			}
		}
		return out
	})
	return found
}
