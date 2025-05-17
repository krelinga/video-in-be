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
	toWrite := [][]byte{x, []byte("\n")}
	for _, b := range toWrite {
		_, err = writer.Write(b)
		if err != nil {
			return err
		}
	}
	return nil
}
