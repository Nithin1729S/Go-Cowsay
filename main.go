package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

// Constants for balloon borders
const (
	topLeft     = "/"
	topRight    = "\\"
	bottomLeft  = "\\"
	bottomRight = "/"
	side        = "|"
	leftCorner  = "<"
	rightCorner = ">"
)

// buildBalloon creates the speech balloon based on the input lines and the maximum width.
func buildBalloon(lines []string, maxWidth int) string {
	var result []string
	top := " " + strings.Repeat("_", maxWidth+2)
	bottom := " " + strings.Repeat("-", maxWidth+2)

	result = append(result, top)

	// Build the lines with borders
	for i, line := range lines {
		switch {
		case len(lines) == 1:
			result = append(result, fmt.Sprintf("%s %s %s", leftCorner, line, rightCorner))
		case i == 0:
			result = append(result, fmt.Sprintf("%s %s %s", topLeft, line, topRight))
		case i == len(lines)-1:
			result = append(result, fmt.Sprintf("%s %s %s", bottomLeft, line, bottomRight))
		default:
			result = append(result, fmt.Sprintf("%s %s %s", side, line, side))
		}
	}

	result = append(result, bottom)
	return strings.Join(result, "\n")
}

// tabsToSpaces converts all tabs to spaces (4 spaces per tab).
func tabsToSpaces(lines []string) []string {
	const tabWidth = 4
	var result []string
	for _, line := range lines {
		result = append(result, strings.ReplaceAll(line, "\t", strings.Repeat(" ", tabWidth)))
	}
	return result
}

// calculateMaxWidth calculates the maximum width of the input lines in terms of runes.
func calculateMaxWidth(lines []string) int {
	maxWidth := 0
	for _, line := range lines {
		if length := utf8.RuneCountInString(line); length > maxWidth {
			maxWidth = length
		}
	}
	return maxWidth
}

// normalizeStringsLength ensures all lines are the same length by padding with spaces.
func normalizeStringsLength(lines []string, maxWidth int) []string {
	var result []string
	for _, line := range lines {
		padding := strings.Repeat(" ", maxWidth-utf8.RuneCountInString(line))
		result = append(result, line+padding)
	}
	return result
}

func main() {
	// Check if the program is being run with piped input
	info, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		os.Exit(1)
	}

	if info.Mode()&os.ModeCharDevice != 0 {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("Usage: fortune | gocowsay")
		return
	}

	// Read input from stdin
	var lines []string
	reader := bufio.NewReader(os.Stdin)

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}
		lines = append(lines, string(line))
	}

	// Define the ASCII art for the cow
	const cow = `
         \  ^__^
          \ (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||
	`

	// Process and normalize input
	lines = tabsToSpaces(lines)
	maxWidth := calculateMaxWidth(lines)
	normalizedLines := normalizeStringsLength(lines, maxWidth)

	// Build and display the speech balloon and cow
	balloon := buildBalloon(normalizedLines, maxWidth)
	fmt.Println(balloon)
	fmt.Println(cow)
}
