package selectlang

import (
	"github.com/mariotoffia/godeviceshadow/model"
	"github.com/mariotoffia/godeviceshadow/model/notifiermodel"
)

// Custom operation value for acknowledge
const (
	MergeOperationAcknowledge model.MergeOperation = 5
)

// LogEntry represents a single log entry in the evaluation context
type LogEntry struct {
	Operation model.MergeOperation
	Path      string
	Value     any
	Keys      map[string]bool // For map values, this contains the keys
}

// EvalContext provides context for evaluating a selection expression
type EvalContext struct {
	// The original operation
	OriginalOp notifiermodel.NotifierOperation

	// The current log entry being evaluated (nil if not evaluating logs)
	CurrentLog *LogEntry

	// Whether we're evaluating in log-entry context
	InLogContext bool
}

// CreateLogEntry creates a LogEntry from various log formats
func CreateLogEntry(operation model.MergeOperation, path string, value any) LogEntry {
	entry := LogEntry{
		Operation: operation,
		Path:      path,
		Value:     value,
		Keys:      make(map[string]bool),
	}

	// If value is a map, extract its keys
	if valueMap, ok := value.(map[string]any); ok {
		for k := range valueMap {
			entry.Keys[k] = true
		}
	}

	return entry
}
