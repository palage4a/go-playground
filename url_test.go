package main

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestUrlParseTcp(t *testing.T) {
	u := "https://debug.test.asga.privet:5000/asdfas/fadsfas?asdfas=asfddas&asdf=fdasf"

	v, err := url.Parse(u)
	if err != nil {
		t.Errorf("failed to parse uri")
	}
	assert.Equal(t, "debug.test.asga.privet:5000", v.Host)
}

func TestUrlParseUnix(t *testing.T) {
	u := "unix:///tmp/test.sock"

	v, err := url.Parse(u)
	if err != nil {
		t.Errorf("failed to parse uri")
	}

	assert.Equal(t, "/tmp/test.sock", v.Path)
}
