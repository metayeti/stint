package tracker

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func WrapText(text string, width int) []string {
	words := strings.Fields(strings.TrimSpace(text))
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		if len(currentLine)+len(word)+1 > width {
			if currentLine != "" {
				lines = append(lines, currentLine)
				currentLine = word
			} else {
				lines = append(lines, word)
				currentLine = ""
			}
		} else {
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	return lines
}

func WriteLogFile(path string, entries []Entry) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	const descriptionWidth = 35 // max width for description column before wrapping

	writer := bufio.NewWriter(file)

	sep := fmt.Sprintf("+--------------+---------------+-%s-+----------------------+\n", strings.Repeat("-", descriptionWidth))
	rowFmt := "| %-12s | %-13s | %-" + strconv.Itoa(descriptionWidth) + "s | %-20s |\n"

	var totalDuration time.Duration

	writer.WriteString(sep)
	writer.WriteString(fmt.Sprintf(rowFmt, "Date", "Time", "Description", "Duration"))
	writer.WriteString(sep)

	for _, entry := range entries {
		descLines := WrapText(entry.Description, descriptionWidth)

		writer.WriteString(fmt.Sprintf(rowFmt, entry.Date, entry.TimeSpan, descLines[0], FormatCustomDuration(entry.Duration)))

		for i := 1; i < len(descLines); i++ {
			writer.WriteString(fmt.Sprintf(rowFmt, "", "", descLines[i], ""))
		}

		totalDuration += entry.Duration
	}

	writer.WriteString(sep)
	writer.WriteString(fmt.Sprintf(rowFmt, "Total", "", "", FormatCustomDuration(totalDuration)))
	writer.WriteString(sep)

	return writer.Flush()
}