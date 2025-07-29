package config

import (
	"os"
	"strconv"
	"strings"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"

	BoldRed    = "\033[1;31m"
	BoldGreen  = "\033[1;32m"
	BoldYellow = "\033[1;33m"
	BoldBlue   = "\033[1;34m"
	BoldPurple = "\033[1;35m"
	BoldCyan   = "\033[1;36m"
	BoldWhite  = "\033[1;37m"
)

// GetTerminalWidth returns the terminal width, defaulting to 80 if unable to determine
func GetTerminalWidth() int {
	// Try to get terminal width from environment variables
	if width := os.Getenv("COLUMNS"); width != "" {
		if widthInt, err := strconv.Atoi(width); err == nil && widthInt > 0 {
			return widthInt
		}
	}

	// Default width if we can't determine
	return 120
}

// CreateSeparator creates a full-width separator with the given character and color
func CreateSeparator(char string, color string) string {
	width := GetTerminalWidth()
	separator := strings.Repeat(char, width)
	return color + separator + Reset
}

// CreateFancySeparator creates a more visually appealing separator with alternating characters
func CreateFancySeparator(color string) string {
	width := GetTerminalWidth()
	pattern := "═"
	separator := strings.Repeat(pattern, width)
	return color + separator + Reset
}

// CreateDottedSeparator creates a dotted separator
func CreateDottedSeparator(color string) string {
	width := GetTerminalWidth()
	pattern := "─"
	separator := strings.Repeat(pattern, width)
	return color + separator + Reset
}

// CreateEmojiSeparator creates a separator with alternating emojis
func CreateEmojiSeparator(color string, emoji string) string {
	width := GetTerminalWidth()
	// Calculate how many emojis can fit (assuming 2 characters per emoji)
	emojiCount := width / 2
	separator := strings.Repeat(emoji+" ", emojiCount)
	// Trim to exact width
	if len(separator) > width {
		separator = separator[:width]
	}
	return color + separator + Reset
}

// CreateDetectorSeparator creates a distinctive separator for detectors
func CreateDetectorSeparator(detectorName string, color string) string {
	width := GetTerminalWidth()
	prefix := "═ " + detectorName + " "
	suffix := " ═"
	availableWidth := width - len(prefix) - len(suffix)
	if availableWidth > 0 {
		fill := strings.Repeat("═", availableWidth)
		return color + prefix + fill + suffix + Reset
	}
	return color + prefix + suffix + Reset
}
