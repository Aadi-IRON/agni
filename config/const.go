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

// CreateBoxHeader creates a beautiful box-style header for detector names
func CreateBoxHeader(title string, color string) string {
	width := GetTerminalWidth()
	boxWidth := width - 4 // Leave 2 spaces on each side

	// Calculate padding for centered title
	titleLen := len(title)
	padding := (boxWidth - titleLen) / 2
	leftPadding := padding
	rightPadding := boxWidth - titleLen - leftPadding

	// Create the box
	topLine := "╭" + strings.Repeat("─", boxWidth) + "╮"
	titleLine := "│" + strings.Repeat(" ", leftPadding) + title + strings.Repeat(" ", rightPadding) + "│"
	bottomLine := "╰" + strings.Repeat("─", boxWidth) + "╯"

	return color + topLine + "\n" + titleLine + "\n" + bottomLine + Reset
}

// CreateSimpleBoxHeader creates a simpler box header with straight corners
func CreateSimpleBoxHeader(title string, color string) string {
	width := GetTerminalWidth()
	boxWidth := width - 4 // Leave 2 spaces on each side

	// Calculate padding for centered title
	titleLen := len(title)
	padding := (boxWidth - titleLen) / 2
	leftPadding := padding
	rightPadding := boxWidth - titleLen - leftPadding

	// Create the box
	topLine := "┌" + strings.Repeat("─", boxWidth) + "┐"
	titleLine := "│" + strings.Repeat(" ", leftPadding) + title + strings.Repeat(" ", rightPadding) + "│"
	bottomLine := "└" + strings.Repeat("─", boxWidth) + "┘"

	return color + topLine + "\n" + titleLine + "\n" + bottomLine + Reset
}

// CreateFancyBoxHeader creates a fancy box header with double lines
func CreateFancyBoxHeader(title string, color string) string {
	width := GetTerminalWidth()
	boxWidth := width - 4 // Leave 2 spaces on each side

	// Ensure title doesn't exceed box width
	if len(title) > boxWidth-2 {
		title = title[:boxWidth-5] + "..."
	}

	// Calculate padding for centered title
	titleLen := len(title)
	padding := (boxWidth - titleLen) / 2
	leftPadding := padding
	rightPadding := boxWidth - titleLen - leftPadding

	// Create the box with double lines
	topLine := "╔" + strings.Repeat("═", boxWidth) + "╗"
	titleLine := "║" + strings.Repeat(" ", leftPadding) + title + strings.Repeat(" ", rightPadding) + "║"
	bottomLine := "╚" + strings.Repeat("═", boxWidth) + "╝"

	return color + topLine + "\n" + titleLine + "\n" + bottomLine + Reset
}

// CreateCompactBoxHeader creates a compact box header that fits better in terminals
func CreateCompactBoxHeader(title string, color string) string {
	// Use a fixed width for better compatibility
	boxWidth := 60

	// Ensure title doesn't exceed box width
	if len(title) > boxWidth-4 {
		title = title[:boxWidth-7] + "..."
	}

	// Calculate padding for centered title
	titleLen := len(title)
	padding := (boxWidth - titleLen) / 2
	leftPadding := padding
	rightPadding := boxWidth - titleLen - leftPadding

	// Create the box
	topLine := "┌" + strings.Repeat("─", boxWidth) + "┐"
	titleLine := "│" + strings.Repeat(" ", leftPadding) + title + strings.Repeat(" ", rightPadding) + "│"
	bottomLine := "└" + strings.Repeat("─", boxWidth) + "┘"

	return color + topLine + "\n" + titleLine + "\n" + bottomLine + Reset
}
