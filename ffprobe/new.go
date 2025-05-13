package ffprobe

import (
	"bytes"
	"encoding/json"
	"os/exec"
)

func New(path string) (*FFProbe, error) {
	cmd := exec.Command("ffprobe", "-i", path, "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", "-select_streams", "v")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var result Raw
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return nil, err
	}

	return &FFProbe{Raw: &result}, nil
}
