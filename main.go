package main

// reads as a DATALOG.TXT file from my Arduino with datalogger shield
// and returns the number and duration of power outages

// the lines in the file are either the string "REBOOT" (at power on)
// or a timestamp (epoch time, but not necessarily based on the correct wall time)

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func formatSeconds(s int64) string {
	h := s / 3600
	m := (s % 3600) / 60
	return fmt.Sprintf("%02d hrs %02d mins", h, m)
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
			duration := formatSeconds(i - previousTimestamp)
			fmt.Printf("Outage #%02d: %s\n", outages, duration)
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
