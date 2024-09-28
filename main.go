package main

// reads as a DATALOG.TXT file from my Arduino with datalogger shield
// and returns the number and duration of power outages

// the lines in the file are either the string "REBOOT" (at power on)
// or a timestamp (epoch time, but not necessarily based on the correct wall time)

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

// a neatly written function that humanizes a time.Duration
// (from https://gist.github.com/harshavardhana/327e0577c4fed9211f65)
func humanizeDuration(duration time.Duration) string {
	days := int64(duration.Hours() / 24)
	hours := int64(math.Mod(duration.Hours(), 24))
	minutes := int64(math.Mod(duration.Minutes(), 60))
	seconds := int64(math.Mod(duration.Seconds(), 60))

	chunks := []struct {
		singularName string
		amount       int64
	}{
		{"day", days},
		{"hour", hours},
		{"minute", minutes},
		{"second", seconds},
	}

	parts := []string{}

	for _, chunk := range chunks {
		switch chunk.amount {
		case 0:
			continue
		case 1:
			parts = append(parts, fmt.Sprintf("%d %s", chunk.amount, chunk.singularName))
		default:
			parts = append(parts, fmt.Sprintf("%d %ss", chunk.amount, chunk.singularName))
		}
	}

	return strings.Join(parts, " ")
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: powersout FILENAME")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	var reboots, unknowns, outages int
	previousLineReboot := false
	var previousTimestamp int64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		text := scanner.Text()

		if strings.Contains(text, "REBOOT") {
			reboots++
			if reboots == 1 {
				fmt.Println()
				continue // skip the first REBOOT
			}
			previousLineReboot = true // for future reference
			continue
		}

		// otherwise this line should be a timestamp

		i, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			unknowns++ // not a timestamp!?
			continue
		}

		if previousLineReboot {
			outages++

			startTime := time.Unix(i, 0)

			duration := time.Duration(i-previousTimestamp) * time.Second
			fmt.Printf("Outage #%02d at %s: %s\n",
				outages,
				startTime.Format("3:04PM on Mon Jan 2, 2006"),
				humanizeDuration(duration),
			)
			previousLineReboot = false
		}

		previousTimestamp = i
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\nReboots: %d\n", reboots-1)
	fmt.Printf("Unknown Entries: %d\n", unknowns)

}
