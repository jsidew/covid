package errors_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	e "github.com/jsidew/covid/internal/errors"
)

const prefix = "test"

func init() {
	e.Prefix = prefix
}

func TestW(t *testing.T) {
	expected := "hello world"
	err := e.W(errors.New("hello world"))
	assert.EqualError(t, err, prefix+": "+expected, "error")
	assert.EqualError(t, errors.Unwrap(err), expected, "unwrapped")
}

func TestF(t *testing.T) {
	expected := "hello kitty"
	err := e.F("hello %s", "kitty")
	assert.EqualError(t, err, prefix+": "+expected)
}
