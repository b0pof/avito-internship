package http

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/b0pof/avito-internship/internal/usecase"
)

type Handler struct {
	uc usecase.IUsecase
}

func NewHandler(uc usecase.IUsecase) *Handler {
	return &Handler{
		uc: uc,
	}
}

func (h *Handler) InitRouter(r *mux.Router) {
	tenders := r.PathPrefix("/tenders").Subrouter()
	{
		tenders.Handle("", http.HandlerFunc(h.GetTenders)).Methods("GET", "OPTIONS")
		tenders.Handle("/new", http.HandlerFunc(h.CreateTender)).Methods("POST", "OPTIONS")
		tenders.Handle("/my", http.HandlerFunc(h.GetMyTenders)).Methods("GET", "OPTIONS")
		tenders.Handle("/{tenderId}/status", http.HandlerFunc(h.GetTenderStatus)).Methods("GET", "OPTIONS")
		tenders.Handle("/{tenderId}/status", http.HandlerFunc(h.UpdateTenderStatus)).Methods("PUT", "OPTIONS")
		tenders.Handle("/{tenderId}/edit", http.HandlerFunc(h.UpdateTender)).Methods("PATCH", "OPTIONS")
		tenders.Handle("/{tenderId}/rollback/{version}", http.HandlerFunc(h.RollbackTender)).Methods("PUT", "OPTIONS")
	}

	bids := r.PathPrefix("/bids").Subrouter()
	{
		bids.Handle("/new", http.HandlerFunc(h.CreateBid)).Methods("POST", "OPTIONS")
		bids.Handle("/my", http.HandlerFunc(h.GetMyBids)).Methods("GET", "OPTIONS")
		bids.Handle("/{tenderId}/list", http.HandlerFunc(h.GetTenderBids)).Methods("GET", "OPTIONS")
		bids.Handle("/{bidId}/status", http.HandlerFunc(h.GetBidStatus)).Methods("GET", "OPTIONS")
		bids.Handle("/{bidId}/status", http.HandlerFunc(h.UpdateBidStatus)).Methods("PUT", "OPTIONS")
		bids.Handle("/{bidId}/submit_decision", http.HandlerFunc(h.SubmitDecision)).Methods("PUT", "OPTIONS")
		bids.Handle("/{bidId}/edit", http.HandlerFunc(h.UpdateBid)).Methods("PATCH", "OPTIONS")
		bids.Handle("/{bidId}/rollback/{version}", http.HandlerFunc(h.RollbackBid)).Methods("PUT", "OPTIONS")
	}
}
