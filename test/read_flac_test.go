package test

import (
	"frolovo22/flac"
	"testing"
)

func TestReadFlac(t *testing.T) {
	_, err := flac.ReadFile("test/BeeMoved.flac")
	if err != nil {
		t.Error(err)
	}
}
