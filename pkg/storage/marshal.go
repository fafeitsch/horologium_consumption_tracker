package storage

import (
	"bytes"
	"fmt"
	"github.com/fafeitsch/Horologium/pkg/consumption"
	"github.com/goccy/go-yaml"
	"io"
)

//Reads the yaml file provided by the reader and returns a series struct.
//In case of parsing errors, an error is returned.
func LoadFromReader(reader io.Reader) (*consumption.Series, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return nil, fmt.Errorf("could not read reader: %v", err)
	}
	series := Series{}
	err = yaml.Unmarshal(buf.Bytes(), &series)
	if err != nil {
		formatError := yaml.FormatError(err, true, true)
		return nil, fmt.Errorf("unmarshalling yaml failed: " + formatError)
	}
	return series.mapToDomain()
}
