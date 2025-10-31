package main

import (
	"LogGenerator/filter"
	"LogGenerator/segment"
	"flag"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// Parse single log line

// func main() {
// 	logLine := `2025-10-23 15:04:10.001 | DEBUG | auth | host=db01 | request_id=req-hyx6sa-8587 | msg="2FA verification completed"`
// 	entry, err := parser.ParseLogEntry(logLine)
// 	if err != nil {
// 		log.Fatal("Error:", err)
// 	}
// 	fmt.Println("Time:", entry.Time.Format("2006-01-02 15:04:05.000"))
// 	fmt.Println("Level:", entry.Level)
// 	fmt.Println("Component:", entry.Component)
// 	fmt.Println("Host:", entry.Host)
// 	fmt.Println("Request ID:", entry.RequestID)
// 	fmt.Println("Message:", entry.Message)
// }

// Parse all log files

// func main() {
// 	entries, _ := parser.ParseLogFiles("../logs")
// 	for _, entry := range entries {
// 		fmt.Println(entry)
// 	}
// 	fmt.Println(len(entries))
// }

// Parse log segments

// func main() {
// 	logStore, err := segment.ParseLogSegments("../logs")
// 	if err != nil {
// 		log.Fatalf("Error parsing log segments: %v", err)
// 	}
// 	if len(logStore.Segments) == 0 {
// 		fmt.Println("No log segments found.")
// 		return
// 	}
// 	seg := logStore.Segments[0]
// 	fmt.Printf("File Name: %s\n", seg.FileName)
// 	fmt.Printf("Start Time: %v\n", seg.StartTime)
// 	fmt.Printf("End Time: %v\n", seg.EndTime)
// 	fmt.Printf("Number of Log Entries: %d\n", len(seg.LogEntries))
// 	fmt.Println("\nLog Entries")
// 	for _, entry := range seg.LogEntries {
// 		fmt.Printf("[%s] | %s | %s | %s | %s\n", entry.Level, entry.Component, entry.Host, entry.RequestID, entry.Message)
// 	}
// }

// Print Index Data

// func main() {
// 	logStore, err := segment.ParseLogSegments("../logs")
// 	if err != nil {
// 		log.Fatalf("Error parsing log segments: %v", err)
// 	}
// 	if len(logStore.Segments) == 0 {
// 		fmt.Println("No log segments found.")
// 		return
// 	}
// 	seg := logStore.Segments[0]

// 	fmt.Println("\nBy Level:")
// 	for level, indices := range seg.Index.ByLevel {
// 		fmt.Printf("  %s → %v\n", level, indices)
// 	}

// 	fmt.Println("\nBy Component:")
// 	for comp, indices := range seg.Index.ByComponent {
// 		fmt.Printf("  %s → %v\n", comp, indices)
// 	}

// 	fmt.Println("\nBy Host:")
// 	for host, indices := range seg.Index.ByHost {
// 		fmt.Printf("  %s → %v\n", host, indices)
// 	}

// 	fmt.Println("\nBy Request ID:")
// 	for req, indices := range seg.Index.ByReqID {
// 		fmt.Printf("  %s → %v\n", req, indices)
// 	}
// }

// Filter command version

func main() {
	level := flag.String("level", "", "Filter by log level")
	component := flag.String("component", "", "Filter by component")
	host := flag.String("host", "", "Filter by host")
	reqID := flag.String("reqID", "", "Filter by requestID")
	startTimeString := flag.String("start", "", "Filter by start time")
	endTimeString := flag.String("end", "", "Filter by end time")
	flag.Parse()

	logStore, err := segment.ParseLogSegments("../logs")

	if err != nil {
		slog.Error("Failed to parse logs\n")
	}
	split := func(s string) []string {
		if s == "" {
			return nil
		}
		parts := strings.Split(s, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return parts
	}
	levels := split(*level)
	components := split(*component)
	hosts := split(*host)
	reqIDs := split(*reqID)
	var startTime, endTime time.Time
	if *startTimeString != "" {
		startTime, err = time.Parse("2006-01-02 15:04:05", *startTimeString)
		if err != nil {
			slog.Error("Error parsing start time", "error", err)
		}
	}

	if *endTimeString != "" {
		endTime, err = time.Parse("2006-01-02 15:04:05", *endTimeString)
		if err != nil {
			slog.Error("Error parsing end time", "error", err)
		}
	}

	filteredLogs := filter.FilterLogs(logStore, levels, components, hosts, reqIDs, startTime, endTime)
	fmt.Printf("Found %d matching entries\n", len(filteredLogs))
	for _, entry := range filteredLogs {
		fmt.Println(entry.Raw)
	}
}
