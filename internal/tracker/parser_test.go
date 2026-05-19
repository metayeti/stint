package tracker

import (
	"os"
	"testing"
	"time"
)

func TestParseCustomDuration(t *testing.T) {
	// Define our test cases
	tests := []struct {
		name    string
		input   string
		want    time.Duration
		wantErr bool
	}{
		{
			name:    "Standard seconds",
			input:   "45s",
			want:    45 * time.Second,
			wantErr: false,
		},
		{
			name:    "Hours and minutes",
			input:   "1h30m",
			want:    90 * time.Minute,
			wantErr: false,
		},
		{
			name:    "Full standard duration",
			input:   "1h23m15s",
			want:    1*time.Hour + 23*time.Minute + 15*time.Second,
			wantErr: false,
		},
		{
			name:    "Custom days format",
			input:   "1d4h5m10s",
			want:    24*time.Hour + 4*time.Hour + 5*time.Minute + 10*time.Second,
			wantErr: false,
		},
		{
			name:    "Just days",
			input:   "2d",
			want:    48 * time.Hour,
			wantErr: false,
		},
		{
			name:    "Invalid format string",
			input:   "not-a-duration",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		// t.Run isolates each test case so one failure doesn't stop the rest
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCustomDuration(tt.input)
			
			// Check if we expected an error but didn't get one, or vice-versa
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseCustomDuration() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			// Verify the parsed duration matches our expected duration exactly
			if got != tt.want {
				t.Errorf("ParseCustomDuration() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatCustomDuration(t *testing.T) {
	tests := []struct {
		name  string
		input time.Duration
		want  string
	}{
		{
			name:  "Seconds only",
			input: 42 * time.Second,
			want:  "42s",
		},
		{
			name:  "Minutes and seconds",
			input: 12 * time.Minute + 5 * time.Second,
			want:  "12m5s",
		},
		{
			name:  "Crosses into days",
			input: 26 * time.Hour + 15 * time.Minute,
			want:  "1d2h15m0s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatCustomDuration(tt.input)
			if got != tt.want {
				t.Errorf("FormatCustomDuration() got = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestParseLogFile(t *testing.T) {
	// 1. Define a sample ASCII table layout representing a realistic file history.
	// This table includes valid rows, a wrapped multi-line description, and decorative boundary lines.
	mockTableContent := `+--------------+---------------+-------------------------------------+----------------------+
| Date         | Time          | Description                         | Duration             |
+--------------+---------------+-------------------------------------+----------------------+
| 2026-05-19   | 10:00 - 11:30 | Setup initial repository design     | 1h30m0s              |
| 2026-05-19   | 14:00 - 15:45 | Fixed authorization edge-case bugs  | 1h45m12s             |
|              |               | and refactored middleware checks    |                      |
| 2026-05-20   | 09:15 - 09:45 | Documentation cleanups              | 30m0s                |
+--------------+---------------+-------------------------------------+----------------------+
| Total        |               |                                     | 3h45m12s             |
+--------------+---------------+-------------------------------------+----------------------+
`

	// 2. Create a temporary file on the filesystem for our test run
	tmpFile, err := os.CreateTemp("", "stint_test_log_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary test file: %v", err)
	}
	// Defers ensure clean-up happens even if a test failure throws a panic mid-run
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// 3. Write our mock table into the temporary file
	if _, err := tmpFile.WriteString(mockTableContent); err != nil {
		t.Fatalf("Failed to populate mock table data: %v", err)
	}
	tmpFile.Close() // Close it to ensure data flushes fully to disk before reading

	// 4. Run your parser against the temp file path
	entries := ParseLogFile(tmpFile.Name())

	// 5. Assert the results match expectations exactly
	expectedLength := 3
	if len(entries) != expectedLength {
		t.Fatalf("ParseLogFile() extracted %d entries, want %d", len(entries), expectedLength)
	}

	// Test Case A: Verify standard one-line row details
	if entries[0].Date != "2026-05-19" || entries[0].Duration != (1*time.Hour + 30*time.Minute) {
		t.Errorf("First entry structural mismatch. Got: %+v", entries[0])
	}

	// Test Case B: Verify multi-line string stitching logic worked smoothly
	expectedStitchedDesc := "Fixed authorization edge-case bugs and refactored middleware checks"
	if entries[1].Description != expectedStitchedDesc {
		t.Errorf("Description unwrapping failed.\nGot:  %q\nWant: %q", entries[1].Description, expectedStitchedDesc)
	}

	// Test Case C: Verify trailing tracking records are safe and untouched by the wrapper loops
	if entries[2].Date != "2026-05-20" || entries[2].Duration != (30 * time.Minute) {
		t.Errorf("Final log record parsed incorrectly. Got: %+v", entries[2])
	}
}