package services

import (
	"testing"
	"time"
)

// TestErrorLogger tests basic error logging functionality
func TestErrorLogger(t *testing.T) {
	logger := NewErrorLogger(100)

	t.Run("LogError", func(t *testing.T) {
		entry := logger.LogError(
			CategoryAmountValidation,
			"Amount cannot be negative",
			map[string]string{"amount": "-100"},
		)

		if entry.Level != ErrorLevelError {
			t.Errorf("expected ErrorLevelError, got %s", entry.Level)
		}
		if entry.Category != CategoryAmountValidation {
			t.Errorf("expected CategoryAmountValidation, got %s", entry.Category)
		}
		if entry.Details["amount"] != "-100" {
			t.Errorf("expected amount -100, got %s", entry.Details["amount"])
		}
	})

	t.Run("LogWarning", func(t *testing.T) {
		entry := logger.LogWarning(
			CategoryAmountValidation,
			"Amount is very large",
			map[string]string{"amount": "999999"},
		)

		if entry.Level != ErrorLevelWarning {
			t.Errorf("expected ErrorLevelWarning, got %s", entry.Level)
		}
	})

	t.Run("LogInfo", func(t *testing.T) {
		entry := logger.LogInfo(
			CategoryDeductionValidation,
			"Deduction created successfully",
			nil,
		)

		if entry.Level != ErrorLevelInfo {
			t.Errorf("expected ErrorLevelInfo, got %s", entry.Level)
		}
	})

	t.Run("GetErrorStats", func(t *testing.T) {
		stats := logger.GetErrorStats()
		if stats[CategoryAmountValidation] < 1 {
			t.Errorf("expected at least 1 error for CategoryAmountValidation")
		}
	})

	t.Run("GetLogs", func(t *testing.T) {
		logs := logger.GetLogs(nil, nil, 100)
		if len(logs) < 3 {
			t.Errorf("expected at least 3 logs, got %d", len(logs))
		}
	})

	t.Run("GetLogsByCategory", func(t *testing.T) {
		category := CategoryAmountValidation
		logs := logger.GetLogs(&category, nil, 100)
		if len(logs) < 2 {
			t.Errorf("expected at least 2 logs for CategoryAmountValidation, got %d", len(logs))
		}
		for _, log := range logs {
			if log.Category != CategoryAmountValidation {
				t.Errorf("expected all logs to be CategoryAmountValidation")
			}
		}
	})

	t.Run("GetLogsByLevel", func(t *testing.T) {
		level := ErrorLevelWarning
		logs := logger.GetLogs(nil, &level, 100)
		if len(logs) < 1 {
			t.Errorf("expected at least 1 warning log, got %d", len(logs))
		}
		for _, log := range logs {
			if log.Level != ErrorLevelWarning {
				t.Errorf("expected all logs to be ErrorLevelWarning")
			}
		}
	})
}

// TestAlertThreshold tests alert threshold functionality
func TestAlertThreshold(t *testing.T) {
	logger := NewErrorLogger(100)

	alertTriggered := false
	var lastAlert AlertEvent

	logger.OnAlert(func(alert AlertEvent) {
		alertTriggered = true
		lastAlert = alert
	})

	// Register threshold: 3 errors in CategoryDeductionValidation trigger alert
	threshold := AlertThreshold{
		Category:   CategoryDeductionValidation,
		ErrorLimit: 3,
		TimeWindow: 1 * time.Minute,
		AlertLevel: ErrorLevelWarning,
	}
	logger.RegisterAlertThreshold(threshold)

	// Log 3 errors to trigger the alert
	for i := 0; i < 3; i++ {
		logger.LogError(
			CategoryDeductionValidation,
			"Deduction validation failed",
			map[string]string{"attempt": string(rune(i + 1))},
		)
	}

	time.Sleep(100 * time.Millisecond)

	if !alertTriggered {
		t.Errorf("expected alert to be triggered after 3 errors")
	}
	if lastAlert.Category != CategoryDeductionValidation {
		t.Errorf("expected alert category to be CategoryDeductionValidation")
	}
	if lastAlert.ErrorCount != 3 {
		t.Errorf("expected error count 3, got %d", lastAlert.ErrorCount)
	}
}

// TestClearLogs tests clearing logs
func TestClearLogs(t *testing.T) {
	logger := NewErrorLogger(100)

	logger.LogError(CategoryAmountValidation, "Error 1", nil)
	logger.LogError(CategoryAmountValidation, "Error 2", nil)

	logs := logger.GetLogs(nil, nil, 100)
	if len(logs) < 2 {
		t.Errorf("expected at least 2 logs before clear")
	}

	logger.ClearLogs()

	logs = logger.GetLogs(nil, nil, 100)
	if len(logs) != 0 {
		t.Errorf("expected 0 logs after clear, got %d", len(logs))
	}

	stats := logger.GetErrorStats()
	if len(stats) != 0 {
		t.Errorf("expected empty stats after clear")
	}
}

// TestMaxLogsLimit tests that logs don't exceed max size
func TestMaxLogsLimit(t *testing.T) {
	maxLogs := 5
	logger := NewErrorLogger(maxLogs)

	// Add more logs than the max
	for i := 0; i < 10; i++ {
		logger.LogError(
			CategoryAmountValidation,
			"Error",
			map[string]string{"index": "value"},
		)
	}

	logs := logger.GetLogs(nil, nil, 100)
	if len(logs) > maxLogs {
		t.Errorf("expected at most %d logs, got %d", maxLogs, len(logs))
	}
}
