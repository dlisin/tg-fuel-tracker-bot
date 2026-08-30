package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAsPointer(t *testing.T) {
	f := func(value any) {
		t.Helper()

		actual := AsPointer(value)

		require.NotNil(t, actual)
		assert.Equal(t, value, *actual)
	}

	f("value")
	f(123)
	f(123.456)
}
