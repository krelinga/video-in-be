package thumbs

import (
	"errors"
	"fmt"

	"github.com/krelinga/video-in-be/state"
)

var (
	ErrQueueFull = errors.New("queue is full")
	ErrUnknownDisc = errors.New("unknown disc")
	ErrDiscState = errors.New("disc state error")
	ErrUnknownProject = errors.New("unknown project")
)

type Queue struct {
	backlog chan *disc
}

func (q *Queue) AddDisc(project, discName string) error {
	if err := setState(project, discName, state.ThumbStateNone, state.ThumbStateWaiting); err != nil {
		return err
	}
	// Needs to return immediately, or maybe panic if the queue is full, cannot block.
	select {
	case q.backlog <- &disc{Project: project, Disc: discName}:
		return nil
	default:
		return ErrQueueFull
	}
}

func NewQueue(length int) *Queue {
	if length <= 0 {
		panic("length must be greater than 0")
	}
	out := &Queue{
		backlog: make(chan *disc, length),
	}
	go generateThumbs(out.backlog)
	return out
}

func setState(project, discName string, oldState, newState state.ThumbState) error {
	var err error
	found := state.ProjectModify(project, func(p *state.Project) {
		disc := p.FindDiscByName(discName)
		if disc == nil {
			err = fmt.Errorf("%w %s for project named %s", ErrUnknownDisc, discName, project)
			return
		}
		if disc.ThumbState != oldState {
			err = fmt.Errorf("%w for project %s disc %s state %s expected state %s", ErrDiscState, project, discName, disc.ThumbState, oldState)
			return
		}
		disc.ThumbState = newState
	})
	if !found {
		err = fmt.Errorf("%w %s", ErrUnknownProject, project)
	}
	return err
}

func trySetError(project, discName string) {
	state.ProjectModify(project, func(p *state.Project) {
		disc := p.FindDiscByName(discName)
		if disc != nil {
			disc.ThumbState = state.ThumbStateError
		}
	})
}