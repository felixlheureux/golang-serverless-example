package httperror

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsErrors(t *testing.T) {
	err1 := CoreUnknownError(errors.New("test my error"))
	assert.True(t, errors.Is(err1, CoreUnknownError(nil)))
}
