package superscript_test

import (
	"caniuse/pkg/superscript"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItoa(t *testing.T) {
	tt := []struct {
		name     string
		input    int
		expected string
	}{
		{
			name:     "lesser than 9",
			input:    7,
			expected: "⁷",
		},
		{
			name:     "greater than 9",
			input:    1234,
			expected: "¹²³⁴",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			expected := superscript.Itoa(tc.input)

			assert.Equal(t, tc.expected, expected)
		})
	}
}
