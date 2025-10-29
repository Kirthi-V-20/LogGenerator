package indexer

import "LogGenerator/model"

func BuildSegmentIndex(LogEntries []model.LogEntry) model.SegmentIndex {
	Index := model.SegmentIndex{
		ByLevel:     make(map[string][]int),
		ByComponent: make(map[string][]int),
		ByHost:      make(map[string][]int),
		ByReqID:     make(map[string][]int),
	}
	for idx, LogEntry := range LogEntries {
		Index.ByLevel[string(LogEntry.Level)] = append(Index.ByLevel[string(LogEntry.Level)], idx)
		Index.ByComponent[LogEntry.Component] = append(Index.ByComponent[LogEntry.Component], idx)
		Index.ByHost[LogEntry.Host] = append(Index.ByHost[LogEntry.Host], idx)
		Index.ByReqID[LogEntry.RequestID] = append(Index.ByReqID[LogEntry.RequestID], idx)
	}
	return Index
}
