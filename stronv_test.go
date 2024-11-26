package main

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStrconvFmt(t *testing.T) {
	a := strconv.FormatInt(time.Now().UnixNano(), 10)
	b := fmt.Sprintf("%d", time.Now().UnixNano())

	assert.Equal(t, a, b)
}
