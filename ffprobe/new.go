package ffprobe

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
)

var (
	ErrFFProbe = errors.New("ffprobe error")
	ErrJson    = errors.New("could not decode ffprobe output as json")
)

func New(path string) (*FFProbe, error) {
	cmd := exec.Command("ffprobe", "-i", path, "-print_format", "json", "-show_format", "-show_streams", "-select_streams", "v")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%w: %w: %s", ErrFFProbe, err, stderr.String())
	}

	var result Raw
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrJson, err)
	}

	return &FFProbe{Raw: &result}, nil
}
