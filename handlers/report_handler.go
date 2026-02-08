package handlers

import (
	"encoding/json"
	"kasir-api/repositories"
	"net/http"
	"time"
)

// ReportHandler mendefinisikan dependency ke repository
type ReportHandler struct {
	repo *repositories.TransactionRepository
}

// NewReportHandler adalah fungsi 'constructor' yang dipanggil di main.go
func NewReportHandler(repo *repositories.TransactionRepository) *ReportHandler {
	return &ReportHandler{repo: repo}
}

// GetDailyReport menangani request GET /api/report/hari-ini
func (h *ReportHandler) GetDailyReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"message": "Method not allowed"})
		return
	}

	// Ambil tanggal hari ini
	today := time.Now().Format("2006-01-02")

	// Panggil fungsi repository yang sudah kita buat sebelumnya
	report, err := h.repo.GetSalesReport(today, today)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// GetReportByRange menangani optional challenge GET /api/report?start_date=...
func (h *ReportHandler) GetReportByRange(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	startDate := query.Get("start_date")
	endDate := query.Get("end_date")

	if startDate == "" || endDate == "" {
		startDate = time.Now().Format("2006-01-02")
		endDate = startDate
	}

	report, err := h.repo.GetSalesReport(startDate, endDate)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}