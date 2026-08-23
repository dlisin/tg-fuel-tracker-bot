package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRegNumber(t *testing.T) {
	f := func(input string, expected RegNumber) {
		t.Helper()

		actual, err := ParseRegNumber(input)

		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	}

	f("A123BC77", RegNumber("A123BC77"))
	f("A123BC777", RegNumber("A123BC777"))
	f("Е347КЕ164", RegNumber("E347KE164"))

	f("a123bc77", RegNumber("A123BC77"))
	f("a123bc777", RegNumber("A123BC777"))

	f("А123ВС77", RegNumber("A123BC77"))
	f("А123ВС777", RegNumber("A123BC777"))

	f("а123вс77", RegNumber("A123BC77"))
	f("а123вс777", RegNumber("A123BC777"))

	f("A123ВС77", RegNumber("A123BC77"))
	f("А123BC77", RegNumber("A123BC77"))

	f("A 123 BC 77", RegNumber("A123BC77"))
	f("А 123 ВС 77", RegNumber("A123BC77"))
	f("а 123 вс 77", RegNumber("A123BC77"))
}

func TestParseRegNumberInvalid(t *testing.T) {
	f := func(input string) {
		t.Helper()

		actual, err := ParseRegNumber(input)

		require.Error(t, err)
		assert.Empty(t, actual)
	}

	f("")
	f("123456")
	f("A123BC")
	f("A123B77")
	f("A1234BC77")
	f("A123@C77")
	f("ABC12377")
	f("A12BC77")
}
