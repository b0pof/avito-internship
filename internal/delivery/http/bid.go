package http

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/b0pof/avito-internship/internal/delivery/dto"
	"github.com/b0pof/avito-internship/internal/model"
	"github.com/b0pof/avito-internship/internal/pkg/helper"
	"github.com/b0pof/avito-internship/internal/repository"
	"github.com/b0pof/avito-internship/internal/usecase"
)

func (h *Handler) CreateBid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bid, err := helper.ParseBidFromBody(r)
	if err != nil {
		helper.Respond(ctx, w, 400, dto.NewErrResponse(err))
		return
	}
	createdTender, err := h.uc.CreateBid(ctx, bid)
	if err != nil {
		var status = 500
		switch {
		case errors.Is(err, model.ErrUserNotFound):
			status = 401
		case errors.Is(err, model.ErrNoRights) || errors.Is(err, model.ErrNoOrganizationFound):
			status = 403
		case errors.Is(err, model.ErrTenderNotFound):
			status = 404
		case errors.Is(err, model.ErrInternal):
			status = 500
		}
		helper.Respond(ctx, w, status, dto.NewErrResponse(err))
		return
	}
	helper.Respond(r.Context(), w, 200, createdTender)
}

func (h *Handler) GetMyBids(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limit, offset, err := helper.ParseLimitOffset(r)
	if err != nil {
		helper.Respond(ctx, w, 400, dto.NewErrResponse(err))
		return
	}
	username := helper.ParseUsername(r)
	tenders, err := h.uc.GetMyBids(ctx, usecase.GetMyBidsInput{
		Limit:    limit,
		Offset:   offset,
		Username: username,
	})
	if err != nil {
		var status = 500
		switch {
		case errors.Is(err, model.ErrUserNotFound):
			status = 401
		case errors.Is(err, model.ErrInternal):
			status = 500
		}
		helper.Respond(ctx, w, status, dto.NewErrResponse(err))
		return
	}
	helper.Respond(r.Context(), w, 200, tenders)
}

func (h *Handler) GetTenderBids(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limit, offset, err := helper.ParseLimitOffset(r)
	if err != nil {
		helper.Respond(ctx, w, 400, dto.NewErrResponse(err))
		return
	}
	username := helper.ParseUsername(r)
	if username == "" {
		helper.Respond(ctx, w, 400, dto.NewErrResponse(model.ErrInvalidQueryParam))
		return
	}
	tenderID := helper.ParseTenderID(r)
	tenders, err := h.uc.GetTenderBids(ctx, repository.GetTenderBidsInput{
		Limit:    limit,
		Offset:   offset,
		Username: username,
		TenderID: tenderID,
	})
	if err != nil {
		var status = 500
		switch {
		case errors.Is(err, model.ErrUserNotFound):
			status = 401
		case errors.Is(err, model.ErrNoRights):
			status = 403
		case errors.Is(err, model.ErrNoBidsFound) || errors.Is(err, model.ErrTenderNotFound):
			status = 404
		case errors.Is(err, model.ErrInternal):
			status = 500
		}
		helper.Respond(ctx, w, status, dto.NewErrResponse(err))
		return
	}
	helper.Respond(r.Context(), w, 200, tenders)
}

func (h *Handler) GetBidStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bidID := helper.ParseBidID(r)
	username := helper.ParseUsername(r)
	stat, err := h.uc.GetBidStatus(ctx, usecase.GetBidStatusInput{
		BidID:    bidID,
		Username: username,
	})
	if err != nil {
		var status = 500
		switch {
		case errors.Is(err, model.ErrUserNotFound):
			status = 401
		case errors.Is(err, model.ErrNoRights):
			status = 403
		case errors.Is(err, model.ErrNoBidFound):
			status = 404
		}
		helper.Respond(ctx, w, status, dto.NewErrResponse(err))
		return
	}
	helper.Respond(r.Context(), w, 200, stat)
}

func (h *Handler) UpdateBidStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bidID := helper.ParseBidID(r)
	username := helper.ParseUsername(r)
	stat := helper.ParseStatus(r)
	bid, err := h.uc.UpdateBidStatus(ctx, usecase.UpdateBidStatusInput{
		BidID:    bidID,
		Status:   stat,
		Username: username,
	})
	if err != nil {
		var status = 500
		switch {
		case errors.Is(err, model.ErrInvalidAttributeValue):
			status = 400
		case errors.Is(err, model.ErrUserNotFound):
			status = 401
		case errors.Is(err, model.ErrNoRights):
			status = 403
		case errors.Is(err, model.ErrNoBidFound):
			status = 404
		}
		helper.Respond(ctx, w, status, dto.NewErrResponse(err))
		return
	}
	helper.Respond(r.Context(), w, 200, bid)
}

func (h *Handler) SubmitDecision(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bidID := helper.ParseBidID(r)
	username := helper.ParseUsername(r)
	decision := helper.ParseDecision(r)
	bid, err := h.uc.SubmitDecision(ctx, usecase.SubmitDecisionInput{
		BidID:    bidID,
		Username: username,
		Decision: decision,
	})
	if err != nil {
		var status = 500
		switch {
		case errors.Is(err, model.ErrWrongDecision):
			status = 400
		case errors.Is(err, model.ErrUserNotFound):
			status = 401
		case errors.Is(err, model.ErrNoRights):
			status = 403
		case errors.Is(err, model.ErrNoBidFound) || errors.Is(err, model.ErrTenderNotFound):
			status = 404
		}
		helper.Respond(ctx, w, status, dto.NewErrResponse(err))
		return
	}
	helper.Respond(r.Context(), w, 200, bid)
}

func (h *Handler) UpdateBid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bidID := helper.ParseBidID(r)
	username := helper.ParseUsername(r)
	info, err := helper.ParseUpdateBidInfo(r)
	if err != nil {
		helper.Respond(ctx, w, 400, dto.NewErrResponse(err))
		return
	}
	updBid, err := h.uc.UpdateBid(ctx, usecase.UpdateBidInput{
		BidID:       bidID,
		Username:    username,
		Name:        info.Name,
		Description: info.Description,
	})
	if err != nil {
		var status = 500
		switch {
		case errors.Is(err, model.ErrInvalidAttributeValue):
			status = 400
		case errors.Is(err, model.ErrUserNotFound):
			status = 401
		case errors.Is(err, model.ErrNoRights):
			status = 403
		case errors.Is(err, model.ErrNoBidFound):
			status = 404
		}
		helper.Respond(ctx, w, status, dto.NewErrResponse(err))
		return
	}
	helper.Respond(r.Context(), w, 200, updBid)
}

func (h *Handler) RollbackBid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bidID := helper.ParseBidID(r)
	username := helper.ParseUsername(r)
	version, err := helper.ParseVersion(r)
	if err != nil {
		helper.Respond(ctx, w, 400, dto.NewErrResponse(err))
		return
	}
	updBid, err := h.uc.RollbackBid(ctx, usecase.RollbackBidInput{
		BidID:    bidID,
		Username: username,
		Version:  version,
	})
	if err != nil {
		var status = 500
		switch {
		case errors.Is(err, model.ErrInvalidAttributeValue):
			status = 400
		case errors.Is(err, model.ErrUserNotFound):
			status = 401
		case errors.Is(err, model.ErrNoRights):
			status = 403
		case errors.Is(err, model.ErrNoBidFound) || errors.Is(err, model.ErrNoSuchVersion):
			status = 404
		}
		helper.Respond(ctx, w, status, dto.NewErrResponse(err))
		return
	}
	helper.Respond(r.Context(), w, 200, updBid)
}
