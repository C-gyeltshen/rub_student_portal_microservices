package services

import (
	"fmt"
	"sync"
	"time"
)

// ErrorLevel represents the severity of an error
type ErrorLevel string

const (
	ErrorLevelDebug   ErrorLevel = "DEBUG"
	ErrorLevelInfo    ErrorLevel = "INFO"
	ErrorLevelWarning ErrorLevel = "WARNING"
	ErrorLevelError   ErrorLevel = "ERROR"
	ErrorLevelCritical ErrorLevel = "CRITICAL"
)

// ErrorCategory represents the category of validation error
type ErrorCategory string

const (
	CategoryAmountValidation      ErrorCategory = "AMOUNT_VALIDATION"
	CategoryDeductionValidation   ErrorCategory = "DEDUCTION_VALIDATION"
	CategoryStipendValidation     ErrorCategory = "STIPEND_VALIDATION"
	CategoryStudentValidation     ErrorCategory = "STUDENT_VALIDATION"
	CategoryBankingValidation     ErrorCategory = "BANKING_VALIDATION"
	CategoryDatabaseError         ErrorCategory = "DATABASE_ERROR"
	CategoryServiceCommunication  ErrorCategory = "SERVICE_COMMUNICATION"
)

// LogEntry represents a single log entry
type LogEntry struct {
	ID        string
	Timestamp time.Time
	Level     ErrorLevel
	Category  ErrorCategory
	Message   string
	Details   map[string]string
}

// ErrorLogger manages error logging with alerting capabilities
type ErrorLogger struct {
	mu              sync.RWMutex
	logs            []LogEntry
	maxLogs         int
	errorCounts     map[ErrorCategory]int
	alertThresholds map[ErrorCategory]AlertThreshold
	alertCallbacks  []func(AlertEvent)
}

// AlertThreshold defines when to trigger an alert
type AlertThreshold struct {
	Category   ErrorCategory
	ErrorLimit int
	TimeWindow time.Duration
	AlertLevel ErrorLevel
}

// AlertEvent represents an alert event
type AlertEvent struct {
	Timestamp  time.Time
	Category   ErrorCategory
	ErrorCount int
	Threshold  int
	Level      ErrorLevel
	Message    string
}

// NewErrorLogger creates a new error logger
func NewErrorLogger(maxLogs int) *ErrorLogger {
	return &ErrorLogger{
		logs:            make([]LogEntry, 0, maxLogs),
		maxLogs:         maxLogs,
		errorCounts:     make(map[ErrorCategory]int),
		alertThresholds: make(map[ErrorCategory]AlertThreshold),
		alertCallbacks:  make([]func(AlertEvent), 0),
	}
}

// Log adds a log entry
func (el *ErrorLogger) Log(level ErrorLevel, category ErrorCategory, message string, details map[string]string) LogEntry {
	el.mu.Lock()
	defer el.mu.Unlock()

	entry := LogEntry{
		ID:        fmt.Sprintf("log_%d", len(el.logs)),
		Timestamp: time.Now(),
		Level:     level,
		Category:  category,
		Message:   message,
		Details:   details,
	}

	// Add to logs
	el.logs = append(el.logs, entry)
	if len(el.logs) > el.maxLogs {
		el.logs = el.logs[1:]
	}

	// Track error counts
	if level == ErrorLevelError || level == ErrorLevelCritical {
		el.errorCounts[category]++

		// Check alert threshold
		if threshold, exists := el.alertThresholds[category]; exists {
			if el.errorCounts[category] >= threshold.ErrorLimit {
				el.triggerAlert(AlertEvent{
					Timestamp:  time.Now(),
					Category:   category,
					ErrorCount: el.errorCounts[category],
					Threshold:  threshold.ErrorLimit,
					Level:      threshold.AlertLevel,
					Message:    fmt.Sprintf("Alert threshold exceeded for %s: %d errors", category, el.errorCounts[category]),
				})
				// Reset counter after alert
				el.errorCounts[category] = 0
			}
		}
	}

	return entry
}

// LogError logs an error
func (el *ErrorLogger) LogError(category ErrorCategory, message string, details map[string]string) LogEntry {
	return el.Log(ErrorLevelError, category, message, details)
}

// LogWarning logs a warning
func (el *ErrorLogger) LogWarning(category ErrorCategory, message string, details map[string]string) LogEntry {
	return el.Log(ErrorLevelWarning, category, message, details)
}

// LogInfo logs info
func (el *ErrorLogger) LogInfo(category ErrorCategory, message string, details map[string]string) LogEntry {
	return el.Log(ErrorLevelInfo, category, message, details)
}

// RegisterAlertThreshold registers an alert threshold for a category
func (el *ErrorLogger) RegisterAlertThreshold(threshold AlertThreshold) {
	el.mu.Lock()
	defer el.mu.Unlock()
	el.alertThresholds[threshold.Category] = threshold
}

// OnAlert registers a callback for alert events
func (el *ErrorLogger) OnAlert(callback func(AlertEvent)) {
	el.mu.Lock()
	defer el.mu.Unlock()
	el.alertCallbacks = append(el.alertCallbacks, callback)
}

// triggerAlert triggers all registered alert callbacks
func (el *ErrorLogger) triggerAlert(alert AlertEvent) {
	for _, callback := range el.alertCallbacks {
		go callback(alert)
	}
}

// GetLogs returns recent logs, optionally filtered by category
func (el *ErrorLogger) GetLogs(category *ErrorCategory, level *ErrorLevel, limit int) []LogEntry {
	el.mu.RLock()
	defer el.mu.RUnlock()

	result := make([]LogEntry, 0)
	for _, log := range el.logs {
		if category != nil && log.Category != *category {
			continue
		}
		if level != nil && log.Level != *level {
			continue
		}
		result = append(result, log)
	}

	// Return last 'limit' entries
	if len(result) > limit {
		return result[len(result)-limit:]
	}
	return result
}

// GetErrorStats returns error statistics
func (el *ErrorLogger) GetErrorStats() map[ErrorCategory]int {
	el.mu.RLock()
	defer el.mu.RUnlock()

	stats := make(map[ErrorCategory]int)
	for k, v := range el.errorCounts {
		stats[k] = v
	}
	return stats
}

// ClearLogs clears all logs
func (el *ErrorLogger) ClearLogs() {
	el.mu.Lock()
	defer el.mu.Unlock()
	el.logs = make([]LogEntry, 0, el.maxLogs)
	el.errorCounts = make(map[ErrorCategory]int)
}
