package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type configuration struct {
	LogFile   string
	Threshold int
}

type logEntry struct {
	Address      string
	Method       string
	URI          string
	ResponseCode string
}

func process(config *configuration) map[string]int {
	file, err := os.Open(config.LogFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	entries := make(map[string]int)

	for _, entry := range lines {
		parts := strings.Split(entry, " ")
		l := logEntry{
			Address:      parts[0],
			Method:       strings.Replace(parts[5], "\"", "", 1),
			URI:          parts[6],
			ResponseCode: parts[8],
		}

		if l.Method == "POST" && l.ResponseCode == "200" {
			entries[l.Address]++
		}
	}

	return entries
}

func report(config *configuration, entries map[string]int) {
	if len(entries) > 0 {
		for k, v := range entries {
			if v >= config.Threshold {
				fmt.Printf("Brute force login attack from %s. Threshold: %d, Actual: %d\n", k, config.Threshold, v)
			}
		}
	}
}

func main() {
	logFilePtr := flag.String("logfile", "app.log", "Path to logfile")
	thresholdPtr := flag.Int("threshold", 10, "Threshold before blacklisting")
	flag.Parse()

	config := &configuration{LogFile: *logFilePtr, Threshold: *thresholdPtr}

	entries := process(config)
	report(config, entries)
}
