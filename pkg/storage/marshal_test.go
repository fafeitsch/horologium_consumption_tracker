package storage

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"strings"
	"testing"
)

func TestLoadFromReader(t *testing.T) {
	file, _ := ioutil.ReadFile("../test-resources/series/powerSeries.yml")
	buffer := bytes.NewBuffer(file)
	got, err := LoadFromReader(buffer)
	require.NoError(t, err, "no error expected")
	assert.Equal(t, "A pseudo power consumption for testing", got.Name, "name not correct")
	require.Equal(t, 3, len(got.PricingPlans))
	assert.Equal(t, "2019", got.PricingPlans[1].Name)
	require.Equal(t, 3, len(got.MeterReadings))
	assert.Equal(t, 1299.23, got.MeterReadings[2].Count)
}

type errReader struct {
}

func (e *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestLoadFromReader_ReaderError(t *testing.T) {
	got, err := LoadFromReader(&errReader{})
	assert.EqualError(t, err, "could not read reader: test error", "error message wrong")
	assert.Nil(t, got, "result should be nil in case of an error")
}

func TestLoadFromReader_YamlError(t *testing.T) {
	reader := strings.NewReader("I'm not { a valid yaml")
	got, err := LoadFromReader(reader)
	assert.EqualError(t, err, "unmarshalling yaml failed: String node doesn't MapNode", "error message wrong")
	assert.Nil(t, got, "result should be nil in case of an error")
}
