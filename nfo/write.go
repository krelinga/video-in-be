package nfo

import (
	"encoding/xml"
	"io"
)

func (m *Movie) Write(writer io.Writer) error {
	x, err := xml.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	_, err = writer.Write(x)
	if err != nil {
		return err
	}
	return nil
}
