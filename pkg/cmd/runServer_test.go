package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_getFlags(t *testing.T) {
	actual := getFlags()
	require.Equal(t, len(actual), 4, "Number of flags incorrect.")
}
