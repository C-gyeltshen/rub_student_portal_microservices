package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"finance_service/database"
	"finance_service/services"
)

// ReportHandler handles report generation requests
type ReportHandler struct {
	reportService *services.ReportService
}

// NewReportHandler creates a new report handler
func NewReportHandler() *ReportHandler {
	return &ReportHandler{
		reportService: services.NewReportService(database.DB),
	}
}

// GetDisbursementReport returns a disbursement summary report
func (rh *ReportHandler) GetDisbursementReport(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate *time.Time

	if startDateStr != "" {
		if t, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			startDate = &t
		}
	}
	if endDateStr != "" {
		if t, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			endDate = &t
		}
	}

	report, err := rh.reportService.GetDisbursementReport(startDate, endDate)
	if err != nil {
		log.Printf("Error generating disbursement report: %v", err)
		http.Error(w, fmt.Sprintf("Failed to generate report: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// GetDeductionReport returns a detailed deduction summary report
func (rh *ReportHandler) GetDeductionReport(w http.ResponseWriter, r *http.Request) {
	reports, err := rh.reportService.GetDeductionReport()
	if err != nil {
		log.Printf("Error generating deduction report: %v", err)
		http.Error(w, fmt.Sprintf("Failed to generate report: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"deductions": reports,
		"count":      len(reports),
	})
}

// GetTransactionReport returns a transaction summary report
func (rh *ReportHandler) GetTransactionReport(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate *time.Time

	if startDateStr != "" {
		if t, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			startDate = &t
		}
	}
	if endDateStr != "" {
		if t, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			endDate = &t
		}
	}

	report, err := rh.reportService.GetTransactionReport(startDate, endDate)
	if err != nil {
		log.Printf("Error generating transaction report: %v", err)
		http.Error(w, fmt.Sprintf("Failed to generate report: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// ExportStipendsCsv exports stipends to CSV format
func (rh *ReportHandler) ExportStipendsCsv(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate *time.Time

	if startDateStr != "" {
		if t, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			startDate = &t
		}
	}
	if endDateStr != "" {
		if t, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			endDate = &t
		}
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=stipends.csv")

	if err := rh.reportService.ExportStipendsToCsv(w, startDate, endDate); err != nil {
		log.Printf("Error exporting stipends to CSV: %v", err)
		http.Error(w, fmt.Sprintf("Failed to export CSV: %v", err), http.StatusInternalServerError)
	}
}

// ExportDeductionsCsv exports deductions to CSV format
func (rh *ReportHandler) ExportDeductionsCsv(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate *time.Time

	if startDateStr != "" {
		if t, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			startDate = &t
		}
	}
	if endDateStr != "" {
		if t, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			endDate = &t
		}
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=deductions.csv")

	if err := rh.reportService.ExportDeductionsToCsv(w, startDate, endDate); err != nil {
		log.Printf("Error exporting deductions to CSV: %v", err)
		http.Error(w, fmt.Sprintf("Failed to export CSV: %v", err), http.StatusInternalServerError)
	}
}

// ExportTransactionsCsv exports transactions to CSV format
func (rh *ReportHandler) ExportTransactionsCsv(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate *time.Time

	if startDateStr != "" {
		if t, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			startDate = &t
		}
	}
	if endDateStr != "" {
		if t, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			endDate = &t
		}
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=transactions.csv")

	if err := rh.reportService.ExportTransactionsToCsv(w, startDate, endDate); err != nil {
		log.Printf("Error exporting transactions to CSV: %v", err)
		http.Error(w, fmt.Sprintf("Failed to export CSV: %v", err), http.StatusInternalServerError)
	}
}
