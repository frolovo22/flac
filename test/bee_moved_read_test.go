package test

import (
	"frolovo22/flac"
	"testing"
)

func TestReadBeeMoved(t *testing.T) {
	_, err := flac.ReadFile("test/BeeMoved.flac")
	if err != nil {
		t.Error(err)
	}
}
