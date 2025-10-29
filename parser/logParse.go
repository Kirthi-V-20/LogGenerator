package parser

import (
	"LogGenerator/model"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func ParseLogEntry(s string) (*model.LogEntry, error) {
	pattern := `^(?P<time>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+)\s+\|\s+(?P<level>[A-Z]+)\s+\|\s+(?P<component>[\w-]+)\s+\|\s+host=(?P<host>[\w-]+)\s+\|\s+request_id=(?P<request_id>[\w-]+)\s+\|\s+msg="(?P<msg>.*)"$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		return nil, fmt.Errorf("Invalid format")
	}
	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if name != "" {
			result[name] = matches[i]
		}
	}
	t, err := time.Parse("2006-01-02 15:04:05.000", result["time"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %v", err)
	}
	entry := model.LogEntry{
		Time:      t,
		Level:     model.LogLevel(result["level"]),
		Component: result["component"],
		Host:      result["host"],
		RequestID: result["request_id"],
		Message:   result["msg"],
		Raw:       matches[0],
	}
	return &entry, nil
}

func ParseLogFiles(s string) ([]model.LogEntry, error) {
	var allEntries []model.LogEntry
	files, err := os.ReadDir(s)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory : %v", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filepath := filepath.Join(s, file.Name())
		f, err := os.Open(filepath)
		if err != nil {
			fmt.Printf("Skipping file %s due to error: %v\n", filepath, err)
			continue
		}
		scanner := bufio.NewScanner(f)
		scanner.Buffer(make([]byte, 0, 1024*1024), 10*1024*1024)
		for scanner.Scan() {
			line := scanner.Text()
			entry, err := ParseLogEntry(line)
			if err == nil {
				allEntries = append(allEntries, *entry)
			}
		}
		err = scanner.Err()
		if err != nil {
			fmt.Printf("Error reading file %s due to error :%v.", filepath, err)
		}
		f.Close()
	}
	return allEntries, nil
}

// func buildSegmentIndex(entries []model.LogEntry) model.SegmentIndex {
// 	index := model.SegmentIndex{
// 		ByLevel:     make(map[string][]int),
// 		ByComponent: make(map[string][]int),
// 		ByHost:      make(map[string][]int),
// 		ByReqID:     make(map[string][]int),
// 	}

// 	for i, entry := range entries {
// 		index.ByLevel[string(entry.Level)] = append(index.ByLevel[string(entry.Level)], i)
// 		index.ByComponent[entry.Component] = append(index.ByComponent[entry.Component], i)
// 		index.ByHost[entry.Host] = append(index.ByHost[entry.Host], i)
// 		index.ByReqID[entry.RequestID] = append(index.ByReqID[entry.RequestID], i)
// 	}

// 	return index
// }
