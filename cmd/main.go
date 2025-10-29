package main

import (
	"LogGenerator/segment"
	"fmt"
	"log"
)

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

// func main() {
// 	entries, _ := parser.ParseLogFiles("../logs")
// 	for _, entry := range entries {
// 		fmt.Println(entry)
// 	}
// 	fmt.Println(len(entries))
// }

// segment

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
// 	fmt.Println("\n Log Entries")

// 	for _, entry := range seg.LogEntries {
// 		fmt.Printf("[%s]  |  %s  |  %s  |  %s  | %s\n", entry.Level, entry.Component, entry.Host, entry.RequestID, entry.Message)
// 	}

// }

func main() {
	logStore, err := segment.ParseLogSegments("../logs")
	if err != nil {
		log.Fatalf("Error parsing log segments: %v", err)
	}
	if len(logStore.Segments) == 0 {
		fmt.Println("No log segments found.")
		return
	}
	seg := logStore.Segments[0]
	fmt.Println("\nBy Level:")
	for level, indices := range seg.Index.ByLevel {
		fmt.Printf("  %s → %v\n", level, indices)
	}
	fmt.Println("\nBy Component:")
	for comp, indices := range seg.Index.ByComponent {
		fmt.Printf("  %s → %v\n", comp, indices)
	}
	fmt.Println("\nBy Host:")
	for host, indices := range seg.Index.ByHost {
		fmt.Printf("  %s → %v\n", host, indices)
	}
	fmt.Println("\nBy Request ID:")
	for req, indices := range seg.Index.ByReqID {
		fmt.Printf("  %s → %v\n", req, indices)
	}
}
