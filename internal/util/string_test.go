package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseInt64(t *testing.T) {
	f := func(input string, expected int64) {
		t.Helper()

		actual, err := ParseInt64(input)

		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	}

	f("0", 0)
	f("123 ", 123)
	f(" -123", -123)
	f(" +123 ", 123)
}

func TestParseInt64Invalid(t *testing.T) {
	f := func(input string) {
		t.Helper()

		actual, err := ParseInt64(input)

		require.Error(t, err)
		assert.Zero(t, actual)
	}

	f("")
	f(" ")
	f("abc")
	f("12.3")
	f("12,3")
	f("1 2")
}

func TestParseFloat64(t *testing.T) {
	f := func(input string, expected float64) {
		t.Helper()

		actual, err := ParseFloat64(input)

		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	}

	f("0", 0)
	f("123 ", 123.0)
	f(" -123", -123.0)
	f(" 123.45 ", 123.45)
	f(" 123,45 ", 123.45)
	f("-123.45", -123.45)
	f("-123,45", -123.45)
}

func TestParseFloat64Invalid(t *testing.T) {
	f := func(input string) {
		t.Helper()

		actual, err := ParseFloat64(input)

		require.Error(t, err)
		assert.Zero(t, actual)
	}

	f("")
	f(" ")
	f("abc")
	f("1 2")
	f("12.3.4")
	f("12,3,4")
}
