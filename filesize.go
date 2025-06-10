// Package filesize provides functionality for parsing human-readable file size
// strings into byte counts. It supports both binary (1024-based) and decimal
// (1000-based) units with various formatting options.
package filesize

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// unit constants for binary (1024-based) calculations
const (
	// byte is the base unit
	Byte int64 = 1

	// binary units (1024-based)
	KiB = Byte * 1024
	MiB = KiB * 1024
	GiB = MiB * 1024
	TiB = GiB * 1024
	PiB = TiB * 1024

	// decimal units (1000-based)
	KB = Byte * 1000
	MB = KB * 1000
	GB = MB * 1000
	TB = GB * 1000
	PB = TB * 1000
)

// unitMap maps unit strings to their byte multipliers
// this supports various formats and case variations
var unitMap = map[string]int64{
	// bytes
	"b":     Byte,
	"byte":  Byte,
	"bytes": Byte,

	// binary units (1024-based) - standard format
	"kib": KiB,
	"mib": MiB,
	"gib": GiB,
	"tib": TiB,
	"pib": PiB,

	// binary units (1024-based) - short format
	"k": KiB,
	"m": MiB,
	"g": GiB,
	"t": TiB,
	"p": PiB,

	// decimal units (1000-based) - standard format
	"kb": KB,
	"mb": MB,
	"gb": GB,
	"tb": TB,
	"pb": PB,
}

// parseRegex matches a number followed by an optional unit
// this regex captures floating point numbers and various unit formats
var parseRegex = regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*([a-zA-Z]*)$`)

// ParseSize converts a human-readable size string to bytes
//
// Supported formats include:
//   - Plain numbers: "1024", "512"
//   - Binary units: "4k", "4K", "4KiB", "10m", "10M", "10MiB"
//   - Decimal units: "4KB", "10MB"
//   - Floating point: "1.5k", "2.5MB"
//
// Binary units use 1024-based calculations (k/K = 1024 bytes)
// Decimal units use 1000-based calculations (KB = 1000 bytes)
func ParseSize(sizeStr string) (int64, error) {
	// trim whitespace from input string
	sizeStr = strings.TrimSpace(sizeStr)

	// handle empty string
	if sizeStr == "" {
		return 0, fmt.Errorf("empty size string")
	}

	// match the input against our parsing regex
	matches := parseRegex.FindStringSubmatch(sizeStr)
	if matches == nil {
		return 0, fmt.Errorf("invalid size format: %s", sizeStr)
	}

	// extract number and unit from regex matches
	numberStr := matches[1]
	unitStr := strings.ToLower(matches[2])

	// parse the numeric portion as a float to handle decimals
	number, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", numberStr)
	}

	// check for negative numbers
	if number < 0 {
		return 0, fmt.Errorf("size cannot be negative: %f", number)
	}

	// handle case where no unit is specified (assume bytes)
	if unitStr == "" {
		return int64(number), nil
	}

	// look up the unit multiplier in our map
	multiplier, exists := unitMap[unitStr]
	if !exists {
		return 0, fmt.Errorf("unknown unit: %s", unitStr)
	}

	// calculate final byte count
	result := number * float64(multiplier)

	// check for overflow by comparing against max int64
	if result > float64(int64(^uint64(0)>>1)) {
		return 0, fmt.Errorf("size too large: %s", sizeStr)
	}

	return int64(result), nil
}

// FormatSize converts a byte count to a human-readable string using binary units
//
// This function automatically selects the most appropriate unit (KiB, MiB, etc.)
// and formats the result to a reasonable number of decimal places.
func FormatSize(bytes int64) string {
	// handle special cases
	if bytes < 0 {
		return "0 B"
	}
	if bytes < KiB {
		return fmt.Sprintf("%d B", bytes)
	}

	// define units in descending order for formatting
	units := []struct {
		name       string
		multiplier int64
	}{
		{"PiB", PiB},
		{"TiB", TiB},
		{"GiB", GiB},
		{"MiB", MiB},
		{"KiB", KiB},
	}

	// find the largest unit that the byte count can be expressed in
	for _, unit := range units {
		if bytes >= unit.multiplier {
			// calculate the value in this unit
			value := float64(bytes) / float64(unit.multiplier)

			// format with appropriate precision
			if value >= 100 {
				// for large values, show no decimal places
				return fmt.Sprintf("%.0f %s", value, unit.name)
			} else if value >= 10 {
				// for medium values, show one decimal place
				return fmt.Sprintf("%.1f %s", value, unit.name)
			} else {
				// for small values, show two decimal places
				return fmt.Sprintf("%.2f %s", value, unit.name)
			}
		}
	}

	// fallback to bytes (should never reach here due to earlier check)
	return fmt.Sprintf("%d B", bytes)
}

// ValidateSize checks if a size string is valid without parsing it
//
// This is useful for validation in CLI flag parsing or configuration
// validation where you want to check format without converting.
func ValidateSize(sizeStr string) error {
	// attempt to parse the size string
	_, err := ParseSize(sizeStr)
	return err
}
