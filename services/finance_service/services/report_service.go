package services

import (
	"encoding/csv"
	"finance_service/models"
	"fmt"
	"io"
	"time"

	"github.com/go-pdf/fpdf"
	"gorm.io/gorm"
)

// ReportService handles report generation
type ReportService struct {
	db *gorm.DB
}

// NewReportService creates a new report service
func NewReportService(db *gorm.DB) *ReportService {
	return &ReportService{db: db}
}

// DisbursementReport contains disbursement summary data
type DisbursementReport struct {
	TotalStipends      int64   `json:"total_stipends"`
	TotalAmount        float64 `json:"total_amount"`
	PendingCount       int64   `json:"pending_count"`
	ProcessedCount     int64   `json:"processed_count"`
	FailedCount        int64   `json:"failed_count"`
	AverageAmount      float64 `json:"average_amount"`
	MinAmount          float64 `json:"min_amount"`
	MaxAmount          float64 `json:"max_amount"`
	GeneratedAt        time.Time `json:"generated_at"`
	ReportPeriod       string  `json:"report_period"`
}

// DeductionReport contains deduction summary data
type DeductionReport struct {
	DeductionType      string  `json:"deduction_type"`
	RuleCount          int64   `json:"rule_count"`
	TotalDeducted      float64 `json:"total_deducted"`
	AverageDeduction   float64 `json:"average_deduction"`
	ApplicableToFull   bool    `json:"applicable_to_full_scholar"`
	ApplicableToSelf   bool    `json:"applicable_to_self_funded"`
}

// TransactionReport contains transaction summary data
type TransactionReport struct {
	TotalTransactions  int64   `json:"total_transactions"`
	SuccessfulCount    int64   `json:"successful_count"`
	PendingCount       int64   `json:"pending_count"`
	FailedCount        int64   `json:"failed_count"`
	TotalAmount        float64 `json:"total_amount"`
	AverageAmount      float64 `json:"average_amount"`
	GeneratedAt        time.Time `json:"generated_at"`
	ReportPeriod       string  `json:"report_period"`
}

// GetDisbursementReport generates a disbursement summary report
func (rs *ReportService) GetDisbursementReport(startDate, endDate *time.Time) (*DisbursementReport, error) {
	var report DisbursementReport

	query := rs.db
	if startDate != nil && endDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
		report.ReportPeriod = fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	}

	// Get counts and amounts
	if err := query.Model(&models.Stipend{}).
		Select("COUNT(*) as total_stipends, COALESCE(SUM(amount), 0) as total_amount, "+
			"COALESCE(AVG(amount), 0) as average_amount, COALESCE(MIN(amount), 0) as min_amount, "+
			"COALESCE(MAX(amount), 0) as max_amount").
		Row().
		Scan(&report.TotalStipends, &report.TotalAmount, &report.AverageAmount, &report.MinAmount, &report.MaxAmount); err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get disbursement stats: %w", err)
	}

	// Get status counts
	var pendingCount, processedCount, failedCount int64
	rs.db.Model(&models.Stipend{}).
		Where("payment_status = ?", "Pending").
		Count(&pendingCount)
	rs.db.Model(&models.Stipend{}).
		Where("payment_status = ?", "Processed").
		Count(&processedCount)
	rs.db.Model(&models.Stipend{}).
		Where("payment_status = ?", "Failed").
		Count(&failedCount)

	report.PendingCount = pendingCount
	report.ProcessedCount = processedCount
	report.FailedCount = failedCount
	report.GeneratedAt = time.Now()

	return &report, nil
}

// GetDeductionReport generates a detailed deduction summary report
func (rs *ReportService) GetDeductionReport() ([]DeductionReport, error) {
	var rules []models.DeductionRule
	if err := rs.db.Find(&rules).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch deduction rules: %w", err)
	}

	var reports []DeductionReport

	for _, rule := range rules {
		var totalDeducted, avgDeduction float64
		var count int64

		// Get deductions for this rule
		rs.db.Model(&models.Deduction{}).
			Where("deduction_rule_id = ?", rule.ID).
			Select("COUNT(*) as count, COALESCE(SUM(amount), 0) as total").
			Row().
			Scan(&count, &totalDeducted)

		if count > 0 {
			avgDeduction = totalDeducted / float64(count)
		}

		reports = append(reports, DeductionReport{
			DeductionType:      rule.DeductionType,
			RuleCount:          1,
			TotalDeducted:      totalDeducted,
			AverageDeduction:   avgDeduction,
			ApplicableToFull:   rule.IsApplicableToFullScholar,
			ApplicableToSelf:   rule.IsApplicableToSelfFunded,
		})
	}

	return reports, nil
}

// GetTransactionReport generates a transaction summary report
func (rs *ReportService) GetTransactionReport(startDate, endDate *time.Time) (*TransactionReport, error) {
	var report TransactionReport

	query := rs.db
	if startDate != nil && endDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
		report.ReportPeriod = fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	}

	// Get overall stats
	if err := query.Model(&models.Transaction{}).
		Select("COUNT(*) as total_transactions, COALESCE(SUM(amount), 0) as total_amount, "+
			"COALESCE(AVG(amount), 0) as average_amount").
		Row().
		Scan(&report.TotalTransactions, &report.TotalAmount, &report.AverageAmount); err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get transaction stats: %w", err)
	}

	// Get status counts
	var successCount, pendingCount, failedCount int64
	rs.db.Model(&models.Transaction{}).
		Where("status = ?", "SUCCESS").
		Count(&successCount)
	rs.db.Model(&models.Transaction{}).
		Where("status = ?", "PENDING").
		Count(&pendingCount)
	rs.db.Model(&models.Transaction{}).
		Where("status = ?", "FAILED").
		Count(&failedCount)

	report.SuccessfulCount = successCount
	report.PendingCount = pendingCount
	report.FailedCount = failedCount
	report.GeneratedAt = time.Now()

	return &report, nil
}

// ExportStipendsToCsv exports stipends to CSV format
func (rs *ReportService) ExportStipendsToCsv(w io.Writer, startDate, endDate *time.Time) error {
	var stipends []models.Stipend

	query := rs.db
	if startDate != nil && endDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	if err := query.Find(&stipends).Error; err != nil {
		return fmt.Errorf("failed to fetch stipends: %w", err)
	}

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "Student ID", "Amount", "Stipend Type", "Payment Status", "Journal Number", "Created At", "Modified At"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data
	for _, stipend := range stipends {
		record := []string{
			stipend.ID.String(),
			stipend.StudentID.String(),
			fmt.Sprintf("%.2f", stipend.Amount),
			stipend.StipendType,
			stipend.PaymentStatus,
			stipend.JournalNumber,
			stipend.CreatedAt.Format(time.RFC3339),
			stipend.ModifiedAt.Format(time.RFC3339),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}

// ExportDeductionsToCsv exports deductions to CSV format
func (rs *ReportService) ExportDeductionsToCsv(w io.Writer, startDate, endDate *time.Time) error {
	var deductions []models.Deduction

	query := rs.db
	if startDate != nil && endDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	if err := query.Find(&deductions).Error; err != nil {
		return fmt.Errorf("failed to fetch deductions: %w", err)
	}

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "Student ID", "Deduction Type", "Amount", "Processing Status", "Deduction Date", "Created At"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data
	for _, deduction := range deductions {
		record := []string{
			deduction.ID.String(),
			deduction.StudentID.String(),
			deduction.DeductionType,
			fmt.Sprintf("%.2f", deduction.Amount),
			deduction.ProcessingStatus,
			deduction.DeductionDate.Format(time.RFC3339),
			deduction.CreatedAt.Format(time.RFC3339),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}

// ExportTransactionsToCsv exports transactions to CSV format
func (rs *ReportService) ExportTransactionsToCsv(w io.Writer, startDate, endDate *time.Time) error {
	var transactions []models.Transaction

	query := rs.db
	if startDate != nil && endDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	if err := query.Find(&transactions).Error; err != nil {
		return fmt.Errorf("failed to fetch transactions: %w", err)
	}

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "Student ID", "Amount", "Status", "Transaction Type", "Destination Account", "Destination Bank", "Reference Number", "Initiated At", "Processed At", "Completed At"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data
	for _, txn := range transactions {
		refNum := ""
		if txn.ReferenceNumber.Valid {
			refNum = txn.ReferenceNumber.String
		}

		record := []string{
			txn.ID.String(),
			txn.StudentID.String(),
			fmt.Sprintf("%.2f", txn.Amount),
			txn.Status,
			txn.TransactionType,
			txn.DestinationAccount,
			txn.DestinationBank,
			refNum,
			txn.InitiatedAt.Format(time.RFC3339),
			txn.ProcessedAt.Format(time.RFC3339),
			txn.CompletedAt.Format(time.RFC3339),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}

// ExportStipendsToPdf exports stipends to PDF format
func (rs *ReportService) ExportStipendsToPdf(w io.Writer, startDate, endDate *time.Time) error {
	var stipends []models.Stipend

	query := rs.db
	if startDate != nil && endDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	if err := query.Find(&stipends).Error; err != nil {
		return fmt.Errorf("failed to fetch stipends: %w", err)
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "STIPENDS REPORT")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(-1)

	if startDate != nil && endDate != nil {
		pdf.Cell(0, 10, fmt.Sprintf("Period: %s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")))
		pdf.Ln(-1)
	}

	var totalAmount float64
	for _, s := range stipends {
		totalAmount += s.Amount
	}

	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 8, "SUMMARY")
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 8, fmt.Sprintf("Total Stipends: %d", len(stipends)))
	pdf.Ln(-1)
	pdf.Cell(0, 8, fmt.Sprintf("Total Amount: Rs. %.2f", totalAmount))
	pdf.Ln(-1)
	if len(stipends) > 0 {
		pdf.Cell(0, 8, fmt.Sprintf("Average Amount: Rs. %.2f", totalAmount/float64(len(stipends))))
		pdf.Ln(-1)
	}

	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(20, 7, "ID")
	pdf.Cell(30, 7, "Student ID")
	pdf.Cell(30, 7, "Amount")
	pdf.Cell(30, 7, "Type")
	pdf.Cell(30, 7, "Status")
	pdf.Cell(20, 7, "Date")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 8)
	for _, stipend := range stipends {
		pdf.Cell(20, 6, stipend.ID.String()[:8])
		pdf.Cell(30, 6, stipend.StudentID.String()[:8])
		pdf.Cell(30, 6, fmt.Sprintf("%.2f", stipend.Amount))
		pdf.Cell(30, 6, stipend.StipendType)
		pdf.Cell(30, 6, stipend.PaymentStatus)
		pdf.Cell(20, 6, stipend.CreatedAt.Format("06-01-02"))
		pdf.Ln(-1)
	}

	pdf.Output(w)
	return nil
}

// ExportDeductionsToPdf exports deductions to PDF format
func (rs *ReportService) ExportDeductionsToPdf(w io.Writer, startDate, endDate *time.Time) error {
	var deductions []models.Deduction

	query := rs.db
	if startDate != nil && endDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	if err := query.Find(&deductions).Error; err != nil {
		return fmt.Errorf("failed to fetch deductions: %w", err)
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "DEDUCTIONS REPORT")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(-1)

	if startDate != nil && endDate != nil {
		pdf.Cell(0, 10, fmt.Sprintf("Period: %s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")))
		pdf.Ln(-1)
	}

	var totalAmount float64
	deductionMap := make(map[string]float64)

	for _, d := range deductions {
		totalAmount += d.Amount
		deductionMap[d.DeductionType] += d.Amount
	}

	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 8, "SUMMARY")
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 8, fmt.Sprintf("Total Deductions: %d", len(deductions)))
	pdf.Ln(-1)
	pdf.Cell(0, 8, fmt.Sprintf("Total Amount: Rs. %.2f", totalAmount))
	pdf.Ln(-1)

	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 8, "Breakdown by Type:")
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 9)

	for deductionType, amount := range deductionMap {
		pdf.Cell(0, 6, fmt.Sprintf("  • %s: Rs. %.2f", deductionType, amount))
		pdf.Ln(-1)
	}

	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(20, 7, "ID")
	pdf.Cell(25, 7, "Student ID")
	pdf.Cell(35, 7, "Deduction Type")
	pdf.Cell(25, 7, "Amount")
	pdf.Cell(20, 7, "Date")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 8)
	for _, deduction := range deductions {
		pdf.Cell(20, 6, deduction.ID.String()[:8])
		pdf.Cell(25, 6, deduction.StudentID.String()[:8])
		pdf.Cell(35, 6, deduction.DeductionType)
		pdf.Cell(25, 6, fmt.Sprintf("%.2f", deduction.Amount))
		pdf.Cell(20, 6, deduction.CreatedAt.Format("06-01-02"))
		pdf.Ln(-1)
	}

	pdf.Output(w)
	return nil
}

// ExportTransactionsToPdf exports transactions to PDF format
func (rs *ReportService) ExportTransactionsToPdf(w io.Writer, startDate, endDate *time.Time) error {
	var transactions []models.Transaction

	query := rs.db
	if startDate != nil && endDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	if err := query.Find(&transactions).Error; err != nil {
		return fmt.Errorf("failed to fetch transactions: %w", err)
	}

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "TRANSACTIONS REPORT")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(-1)

	if startDate != nil && endDate != nil {
		pdf.Cell(0, 10, fmt.Sprintf("Period: %s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")))
		pdf.Ln(-1)
	}

	var totalAmount float64
	successCount, pendingCount, failedCount := 0, 0, 0

	for _, t := range transactions {
		totalAmount += t.Amount
		switch t.Status {
		case "SUCCESS":
			successCount++
		case "PENDING":
			pendingCount++
		case "FAILED":
			failedCount++
		}
	}

	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 8, "SUMMARY")
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 8, fmt.Sprintf("Total Transactions: %d", len(transactions)))
	pdf.Ln(-1)
	pdf.Cell(0, 8, fmt.Sprintf("Successful: %d | Pending: %d | Failed: %d", successCount, pendingCount, failedCount))
	pdf.Ln(-1)
	pdf.Cell(0, 8, fmt.Sprintf("Total Amount: Rs. %.2f", totalAmount))
	pdf.Ln(-1)

	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 8)
	pdf.Cell(15, 7, "ID")
	pdf.Cell(15, 7, "Student")
	pdf.Cell(15, 7, "Amount")
	pdf.Cell(12, 7, "Status")
	pdf.Cell(15, 7, "Type")
	pdf.Cell(25, 7, "Account")
	pdf.Cell(15, 7, "Date")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 7)
	for _, txn := range transactions {
		pdf.Cell(15, 6, txn.ID.String()[:8])
		pdf.Cell(15, 6, txn.StudentID.String()[:8])
		pdf.Cell(15, 6, fmt.Sprintf("%.0f", txn.Amount))
		pdf.Cell(12, 6, txn.Status)
		pdf.Cell(15, 6, txn.TransactionType)
		acctDisplay := txn.DestinationAccount
		if len(acctDisplay) > 12 {
			acctDisplay = acctDisplay[:12]
		}
		pdf.Cell(25, 6, acctDisplay)
		pdf.Cell(15, 6, txn.CreatedAt.Format("06-01-02"))
		pdf.Ln(-1)
	}

	pdf.Output(w)
	return nil
}
