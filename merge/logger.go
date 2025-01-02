package merge

import "time"

func (ml MergeLoggers) NotifyProcessed(
	path string,
	operation MergeOperation,
	oldValue, newValue any,
	oldTimeStamp, newTimeStamp time.Time,
) {
	for _, l := range ml {
		l.Processed(path, operation, oldValue, newValue, oldTimeStamp, newTimeStamp)
	}
}
