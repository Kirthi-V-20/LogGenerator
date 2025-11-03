package segment

import (
	"LogGenerator/indexer"
	"LogGenerator/model"
	"LogGenerator/parser"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func ParseLogSegments(path string) (model.LogStore, error) {
	start := time.Now()
	LogStore := model.LogStore{
		Segments: []model.Segment{},
	}
	files, err := os.ReadDir(path)
	if err != nil {
		return LogStore, fmt.Errorf("failed to read directory : %v", err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		wg.Add(1)
		go func(file os.DirEntry) {
			defer wg.Done()

			filePath := filepath.Join(path, file.Name())
			f, err := os.Open(filePath)
			if err != nil {
				fmt.Printf("Skipping file %s due to error: %v", filePath, err)
				return
			}
			defer f.Close()
			var LogEntries []model.LogEntry
			scanner := bufio.NewScanner(f)
			scanner.Buffer(make([]byte, 0, 1024*1024), 10*1024*1024)

			for scanner.Scan() {
				line := scanner.Text()
				entry, err := parser.ParseLogEntry(line)
				if err == nil {
					LogEntries = append(LogEntries, *entry)
				}
			}
			f.Close()
			if err := scanner.Err(); err != nil {
				fmt.Printf("Error reading file %s: %v\n", filePath, err)
			}

			if len(LogEntries) == 0 {
				return
			}
			index := indexer.BuildSegmentIndex(LogEntries)
			segment := model.Segment{
				FileName:   file.Name(),
				LogEntries: LogEntries,
				StartTime:  LogEntries[0].Time,
				EndTime:    LogEntries[len(LogEntries)-1].Time,
				Index:      index,
			}
			mu.Lock()
			LogStore.Segments = append(LogStore.Segments, segment)
			mu.Unlock()
		}(file)
	}

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Println("Segment parsing took:", elapsed)
	return LogStore, nil

}
