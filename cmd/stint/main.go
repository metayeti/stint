/*       _   _       _
 *   ___| |_(_)_ __ | |_
 *  / __| __| | '_ \| __|
 *  \__ \ |_| | | | | |_
 *  |___/\__|_|_| |_|\__|
 *
 *
 *  stint (noun): A fixed or limited period of time spent doing a particular
 *                job or activity.
 *
 *  ---
 *
 *  This is a very simple command-line time tracking utility for developers.
 *  It tracks the amount of time you spend on tasks.
 *
 *  To use, run: stint "your task"
 *
 *  ---
 *
 *  (c) 2026 Danijel Durakovic
 *  Licensed under the terms of the MIT license.
 *
 */

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/metayeti/stint/internal/tracker"
)

const appName = "stint"
const appVersion = "v0.9.9"

const outputDir = "logs"

func printVersionHeader() {
	fmt.Printf(":: %s %s\n\n", hueShiftedString(appName), appVersion)	
}

func main() {

	printVersionHeader()

	const prefix = orange + ">>" + reset;

	if len(os.Args) < 2 {
		fmt.Println("Usage: stint \"task name\"")
		return
	}

	taskName := strings.TrimSpace(os.Args[1])
	if taskName == "" {
		fmt.Println("Error: Task name cannot be empty.")
	}	

	fileName := fmt.Sprintf("%s.txt", taskName)
	filePath := filepath.Join(outputDir, fileName)

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	var existingEntries []tracker.Entry
	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("%s Task log for \"%s\" exists. New tasks will be appended.\n", prefix, taskName)
		existingEntries = tracker.ParseLogFile(filePath)
	} else {
		fmt.Printf("%s Creating a new task log for \"%s\".\n", prefix, taskName)
	}

	startTime := time.Now()

	fmt.Printf(
		"\n%s Tracking started for \"%s\" at %s\n\n", 
		prefix,
		taskName,
		startTime.Format("15:04:01"),
	)
	fmt.Println("Press [ENTER] to stop timing ...")


	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

	endTime := time.Now()
	elapsed := endTime.Sub(startTime).Round(time.Second)

	fmt.Printf("%s Tracking stopped. Elapsed time: %s\n\n", prefix, elapsed)

	fmt.Print("Describe the task you just did: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	if description == "" {
		description = "n/a"
	}

	currentDate := startTime.Format("2006-01-02")
	timeSpan := fmt.Sprintf("%s - %s", startTime.Format("15:04"), endTime.Format("15:04"))

	existingEntries = append(existingEntries, tracker.Entry{
		Date: currentDate,
		TimeSpan: timeSpan,
		Description: description,
		Duration: elapsed,
	})

	if err := tracker.WriteLogFile(filePath, existingEntries); err != nil {
		fmt.Printf("Error saving data to file: %v\n", err)
	}

	fmt.Printf("\n%s Log updated successfully at: %s\n", prefix, filePath)
}