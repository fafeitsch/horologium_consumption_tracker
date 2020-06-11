package storage

import (
	"bytes"
	"fmt"
	"github.com/fafeitsch/Horologium/pkg/domain"
	"github.com/goccy/go-yaml"
	"io"
)

func LoadFromReader(reader io.Reader) (*domain.Series, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return nil, fmt.Errorf("could not read reader: %v", err)
	}
	series := Series{}
	err = yaml.Unmarshal(buf.Bytes(), &series)
	return nil, err
}
