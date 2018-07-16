package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type configuration struct {
	LogFile string
}

type logEntry struct {
	Address      string
	Method       string
	URI          string
	ResponseCode string
}

func process(config *configuration) {
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

	for _, entry := range lines {
		parts := strings.Split(entry, " ")
		l := logEntry{
			Address:      parts[0],
			Method:       strings.Replace(parts[5], "\"", "", 1),
			URI:          parts[6],
			ResponseCode: parts[8],
		}
		fmt.Println(l)
	}
}

func report() {
	fmt.Println("Complete")
}

func main() {
	logFilePtr := flag.String("logfile", "app.log", "Path to logfile")
	flag.Parse()

	config := &configuration{LogFile: *logFilePtr}

	process(config)
	report()
}
