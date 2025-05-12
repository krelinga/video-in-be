package state

import (
	"sync"

	"github.com/krelinga/video-in-be/env"
)

type Project struct {
	Name string `json:"name"`
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
