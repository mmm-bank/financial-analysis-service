package http

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	. "github.com/mmm-bank/financial-analysis-service/storage"
	"github.com/mmm-bank/infra/middleware"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

type TransactionService struct {
	historyDB       TransactionStorage
	analysisStorage AnalysisStorage
	logger          *zap.Logger
}

func NewTransactionService(historyDB TransactionStorage, analysisCollection AnalysisStorage) *TransactionService {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	return &TransactionService{
		historyDB:       historyDB,
		analysisStorage: analysisCollection,
		logger:          logger,
	}
}

func (t *TransactionService) getHistoryHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uuid.UUID)
	transactions, err := t.historyDB.GetTransactions(userID)
	if err != nil {
		t.logger.Error("Error getting transactions", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		t.logger.Error("Error encoding transactions to JSON", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (t *TransactionService) getExpensesHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uuid.UUID)
	monthYear := chi.URLParam(r, "month_year")
	if monthYear == "" {
		http.Error(w, "Missing month parameter", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	analysis, err := t.analysisStorage.GetExpensesAnalysis(ctx, userID, monthYear)
	if err != nil {
		t.logger.Error("Error getting financial analysis", zap.Error(err))
		http.Error(w, "Error getting financial analysis", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(analysis); err != nil {
		t.logger.Error("Error encoding financial analysis", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (t *TransactionService) getIncomeHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(uuid.UUID)
	monthYear := chi.URLParam(r, "month_year")
	if monthYear == "" {
		http.Error(w, "Missing month parameter", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	analysis, err := t.analysisStorage.GetIncomeAnalysis(ctx, userID, monthYear)
	if err != nil {
		t.logger.Error("Error getting financial analysis", zap.Error(err))
		http.Error(w, "Error getting financial analysis", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(analysis); err != nil {
		t.logger.Error("Error encoding financial analysis", zap.Error(err))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func CreateAndRunServer(t *TransactionService, addr string) error {
	r := chi.NewRouter()
	r.Use(mymiddleware.ExtractPayload)
	r.Route("/finance", func(r chi.Router) {
		r.Get("/history", t.getHistoryHandler)
		r.Route("/analysis", func(r chi.Router) {
			r.Get("/expenses/{month_year}", t.getExpensesHandler)
			r.Get("/income/{month_year}", t.getIncomeHandler)
		})
	})
	return http.ListenAndServe(addr, r)
}
