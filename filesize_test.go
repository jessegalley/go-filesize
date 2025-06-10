package filesize

import (
	"testing"
)

// TestParseSize tests the ParseSize function with various input formats
func TestParseSize(t *testing.T) {
	// define test cases with input and expected output
	testCases := []struct {
		input    string
		expected int64
		hasError bool
	}{
		// basic byte values
		{"0", 0, false},
		{"1", 1, false},
		{"1024", 1024, false},
		{"2048", 2048, false},

		// binary units - short format (1024-based)
		{"1k", 1024, false},
		{"1K", 1024, false},
		{"2k", 2048, false},
		{"1m", 1024 * 1024, false},
		{"1M", 1024 * 1024, false},
		{"1g", 1024 * 1024 * 1024, false},
		{"1G", 1024 * 1024 * 1024, false},

		// binary units - standard format (1024-based)
		{"1KiB", 1024, false},
		{"1MiB", 1024 * 1024, false},
		{"1GiB", 1024 * 1024 * 1024, false},
		{"1TiB", 1024 * 1024 * 1024 * 1024, false},

		// decimal units - standard format (1000-based)
		{"1KB", 1000, false},
		{"1MB", 1000 * 1000, false},
		{"1GB", 1000 * 1000 * 1000, false},
		{"1TB", 1000 * 1000 * 1000 * 1000, false},

		// floating point values
		{"1.5k", int64(1.5 * 1024), false},
		{"2.5KiB", int64(2.5 * 1024), false},
		{"1.5MB", int64(1.5 * 1000 * 1000), false},
		{"0.5m", int64(0.5 * 1024 * 1024), false},

		// explicit byte units
		{"100b", 100, false},
		{"100B", 100, false},
		{"100byte", 100, false},
		{"100bytes", 100, false},

		// whitespace handling
		{" 1k ", 1024, false},
		{"1 k", 1024, false},
		{" 1 KiB ", 1024, false},

		// edge cases that should work
		{"4k", 4 * 1024, false}, // your example from the prompt
		{"4kib", 4 * 1024, false},
		{"4KiB", 4 * 1024, false},

		// error cases - invalid format
		{"", 0, true},
		{"abc", 0, true},
		{"k1", 0, true},
		{"1xy", 0, true},
		{"1.2.3k", 0, true},

		// error cases - negative numbers
		{"-1", 0, true},
		{"-1k", 0, true},

		// error cases - unknown units
		{"1XB", 0, true},
		{"1ZiB", 0, true},
	}

	// run each test case
	for _, tc := range testCases {
		result, err := ParseSize(tc.input)

		// check error expectation
		if tc.hasError {
			if err == nil {
				t.Errorf("ParseSize(%q) expected error but got none", tc.input)
			}
			continue
		}

		// check for unexpected errors
		if err != nil {
			t.Errorf("ParseSize(%q) unexpected error: %v", tc.input, err)
			continue
		}

		// check result value
		if result != tc.expected {
			t.Errorf("ParseSize(%q) = %d, expected %d", tc.input, result, tc.expected)
		}
	}
}

// TestParseSize_SpecificExamples tests the specific examples from the prompt
func TestParseSize_SpecificExamples(t *testing.T) {
	// test the example from the prompt: 4k = 4 KiB = 4*1024
	testCases := []struct {
		input    string
		expected int64
	}{
		{"4k", 4 * 1024},
		{"4K", 4 * 1024},
		{"4kib", 4 * 1024},
		{"4KiB", 4 * 1024},
		{"4KB", 4 * 1000}, // note: decimal unit should be 1000-based
	}

	for _, tc := range testCases {
		result, err := ParseSize(tc.input)
		if err != nil {
			t.Errorf("ParseSize(%q) unexpected error: %v", tc.input, err)
			continue
		}

		if result != tc.expected {
			t.Errorf("ParseSize(%q) = %d, expected %d", tc.input, result, tc.expected)
		}
	}
}

// TestFormatSize tests the FormatSize function
func TestFormatSize(t *testing.T) {
	testCases := []struct {
		input    int64
		expected string
	}{
		// basic byte values
		{0, "0 B"},
		{1, "1 B"},
		{512, "512 B"},
		{1023, "1023 B"},

		// kibibytes
		{1024, "1.00 KiB"},
		{2048, "2.00 KiB"},
		{1536, "1.50 KiB"},
		{10240, "10.0 KiB"},
		{102400, "100 KiB"},

		// mebibytes
		{1024 * 1024, "1.00 MiB"},
		{1024 * 1024 * 2, "2.00 MiB"},
		{1024 * 1024 * 10, "10.0 MiB"},
		{1024 * 1024 * 100, "100 MiB"},

		// gibibytes
		{1024 * 1024 * 1024, "1.00 GiB"},
		{1024 * 1024 * 1024 * 2, "2.00 GiB"},

		// negative values (edge case)
		{-1, "0 B"},
	}

	for _, tc := range testCases {
		result := FormatSize(tc.input)
		if result != tc.expected {
			t.Errorf("FormatSize(%d) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}

// TestValidateSize tests the ValidateSize function
func TestValidateSize(t *testing.T) {
	validSizes := []string{
		"1",
		"1k",
		"1KiB",
		"1.5MB",
		"100b",
		" 1 k ",
	}

	invalidSizes := []string{
		"",
		"abc",
		"1xy",
		"-1k",
		"1ZiB",
	}

	// test valid sizes
	for _, size := range validSizes {
		if err := ValidateSize(size); err != nil {
			t.Errorf("ValidateSize(%q) expected no error but got: %v", size, err)
		}
	}

	// test invalid sizes
	for _, size := range invalidSizes {
		if err := ValidateSize(size); err == nil {
			t.Errorf("ValidateSize(%q) expected error but got none", size)
		}
	}
}

// BenchmarkParseSize benchmarks the ParseSize function
func BenchmarkParseSize(b *testing.B) {
	testCases := []string{
		"1k",
		"1KiB",
		"1.5MB",
		"100",
		"1GiB",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			_, _ = ParseSize(tc)
		}
	}
}

// BenchmarkFormatSize benchmarks the FormatSize function
func BenchmarkFormatSize(b *testing.B) {
	testCases := []int64{
		1024,
		1024 * 1024,
		1024 * 1024 * 1024,
		1000 * 1000,
		1000 * 1000 * 1000,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			_ = FormatSize(tc)
		}
	}
}
