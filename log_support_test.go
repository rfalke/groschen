package groschen

import (
	"testing"
)

func TestFormat(t *testing.T) {
	testCases := []struct {
		value    int
		expected string
	}{
		{1, "1"},
		{12, "12"},
		{123, "123"},
		{1234, "1,234"},
		{12345, "12,345"},
		{123456, "123,456"},
		{1234567, "1,234,567"},
		{12345678, "12,345,678"},
		{123456789, "123,456,789"},
		{1234567890, "1,234,567,890"},
		{12345678901, "12,345,678,901"},
	}

	for _, tc := range testCases {
		str := FormatIntWithThousandSeparator(tc.value, ",")
		if str != tc.expected {
			t.Errorf("value=%d got='%s' expected='%s'", tc.value, str, tc.expected)
		}
	}
}
