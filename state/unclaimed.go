package state

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"connectrpc.com/connect"
	"github.com/krelinga/video-in-be/env"
)

var unclaimedMutex = &sync.RWMutex{}

var (
	ErrUnclaimedDirNotFound = connect.NewError(connect.CodeNotFound, errors.New("unclaimed directory not found"))
	ErrUnclaimedDirMoveFailed  = connect.NewError(connect.CodeDataLoss, errors.New("unclaimed directory move failed"))
	ErrUnknownProject = connect.NewError(connect.CodeNotFound, errors.New("unknown project"))
)

func listUnclaimedDirs() []string {
	entries, err := os.ReadDir(env.UnclaimedDir())
	if err != nil {
		panic(fmt.Sprint("Could not read unclaimed directory: ", err))
	}
	dirs := []string{}
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e.Name())
		}
	}
	return dirs
}

func UnclaimedDiscDirRead(f func( []string)) {
	unclaimedMutex.RLock()
	defer unclaimedMutex.RUnlock()

	f(listUnclaimedDirs())
}

func ProjectAssignDiskDirs(project string, dirs []string) error {
	unclaimedMutex.Lock()
	defer unclaimedMutex.Unlock()

	// Make sure every directory exists
	foundDirs := map[string]struct{}{}
	for _, dir := range listUnclaimedDirs() {
		foundDirs[dir] = struct{}{}
	}
	for _, dir := range dirs {
		if _, ok := foundDirs[dir]; !ok {
			return ErrUnclaimedDirNotFound
		}
	}

	// Make sure the project exists and do the move
	var err error
	ProjectsRead(func(projects []*Project) {
		for _, p := range projects {
			if p.Name == project {
				for _, dir := range dirs {
					mvErr := os.Rename(env.UnclaimedDir()+"/"+dir, env.StateDir()+"/"+project+"/"+dir)
					if mvErr != nil {
						err = ErrUnclaimedDirMoveFailed
						return
					}
				}
				return
			}
		}
		err = ErrUnknownProject
	})

	return err
}