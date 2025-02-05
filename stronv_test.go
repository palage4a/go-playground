package main

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStrconvFmt(t *testing.T) {
	ts := time.Now().UnixNano()
	a := strconv.FormatInt(ts, 10)
	b := fmt.Sprintf("%d", ts)

	assert.Equal(t, a, b)
}
