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

type Project struct {
	Name   string                `json:"name"`
	Thumbs map[string]ThumbState `json:"thumbs,omitempty"`
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
