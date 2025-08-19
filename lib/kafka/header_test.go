package kafka

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaderGet(t *testing.T) {
	testcases := []struct {
		name     string
		header   Header
		key      string
		expected string
	}{
		{
			name:     "key exists",
			header:   Header{"key": "value"},
			key:      "key",
			expected: "value",
		},
		{
			name:     "key does not exist",
			header:   Header{"key": "value"},
			key:      "not-exist",
			expected: "",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.header.Get(tc.key)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestHeaderSet(t *testing.T) {
	testcases := []struct {
		name     string
		header   Header
		key      string
		value    string
		expected Header
	}{
		{
			name:     "set key",
			header:   Header{"key": "value"},
			key:      "new-key",
			value:    "new-value",
			expected: Header{"key": "value", "new-key": "new-value"},
		},
		{
			name:     "update key",
			header:   Header{"key": "value"},
			key:      "key",
			value:    "new-value",
			expected: Header{"key": "new-value"},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.header.Set(tc.key, tc.value)
			assert.Equal(t, tc.expected, tc.header)
		})
	}
}

func TestHeaderDel(t *testing.T) {
	testcases := []struct {
		name     string
		header   Header
		key      string
		expected Header
	}{
		{
			name:     "delete key",
			header:   Header{"key": "value"},
			key:      "key",
			expected: Header{},
		},
		{
			name:     "key does not exist",
			header:   Header{"key": "value"},
			key:      "not-exist",
			expected: Header{"key": "value"},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.header.Del(tc.key)
			assert.Equal(t, tc.expected, tc.header)
		})
	}
}

func TestHeaderKeys(t *testing.T) {
	testcases := []struct {
		name     string
		header   Header
		expected []string
	}{
		{
			name:     "empty header",
			header:   Header{},
			expected: []string{},
		},
		{
			name:     "non-empty header",
			header:   Header{"key": "value", "new-key": "new-value"},
			expected: []string{"key", "new-key"},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.header.Keys()
			assert.ElementsMatch(t, tc.expected, actual)
		})
	}
}

func TestHeaderValues(t *testing.T) {
	testcases := []struct {
		name     string
		header   Header
		expected []string
	}{
		{
			name:     "empty header",
			header:   Header{},
			expected: []string{},
		},
		{
			name:     "non-empty header",
			header:   Header{"key": "value", "new-key": "new-value"},
			expected: []string{"value", "new-value"},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.header.Values()
			assert.ElementsMatch(t, tc.expected, actual)
		})
	}
}
