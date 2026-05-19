package tracker

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

func ParseCustomDuration(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	var days int
	var err error

	if strings.Contains(s, "d") {
		parts := strings.SplitN(s, "d", 2)
		days, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}
		s = parts[1]
	}

	if s == "" {
		return time.Duration(days) * 24 * time.Hour, nil
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, err
	}

	return d + (time.Duration(days) * 24 * time.Hour), nil
}

func FormatCustomDuration(d time.Duration) string {
	secs := int64(d.Seconds())
	days := secs / 86400
	secs %= 86400
	hours := secs / 3600
	secs %= 3600
	minutes := secs / 60
	secs %= 60

	var parts []string
	if days > 0 {
		parts = append(parts, strconv.FormatInt(days, 10)+"d")
	}
	if hours > 0 || days > 0 {
		parts = append(parts, strconv.FormatInt(hours, 10)+"h")
	}
	if minutes > 0 || hours > 0 || days > 0 {
		parts = append(parts, strconv.FormatInt(minutes, 10)+"m")
	}
	parts = append(parts, strconv.FormatInt(secs, 10)+"s")

	return strings.Join(parts, "")
}

func ParseLogFile(path string) []Entry {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var entries []Entry
	scanner := bufio.NewScanner(file)
	var currentEntry *Entry

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "|") {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 6 {
			continue
		}

		datePart := strings.TrimSpace(parts[1])
		timePart := strings.TrimSpace(parts[2])
		descPart := strings.TrimSpace(parts[3])
		durPart := strings.TrimSpace(parts[4])

		// ignore table header rows and the running Total row
		if datePart == "Date" || datePart == "Total" {
			if currentEntry != nil {
				entries = append(entries, *currentEntry)
				currentEntry = nil
			}
			continue
		}

		// handle a new primary tracking record line
		if datePart != "" {
			if currentEntry != nil {
				entries = append(entries, *currentEntry)
			}

			dur, err := ParseCustomDuration(durPart)
			if err != nil {
				currentEntry = nil
				continue
			}

			currentEntry = &Entry{
				Date:        datePart,
				TimeSpan:    timePart,
				Description: descPart,
				Duration:    dur,
			}
		} else if currentEntry != nil && descPart != "" {
			// handle continuation lines for wrapped descriptions.
			if currentEntry.Description != "" {
				currentEntry.Description += " " + descPart
			} else {
				currentEntry.Description = descPart
			}
		}
	}

	// catch any trailing entry left in memory when the scanner finishes
	if currentEntry != nil {
		entries = append(entries, *currentEntry)
	}

	return entries
}