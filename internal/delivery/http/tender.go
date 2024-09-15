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

func (h *Handler) GetTenders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limit, offset, err := helper.ParseLimitOffset(r)
	if err != nil {
		helper.Respond(ctx, w, 400, dto.NewErrResponse(err))
		return
	}
	serviceTypes := helper.ParseServiseTypes(r)

	result, err := h.uc.GetTenders(ctx, repository.GetTendersInput{
		Limit:        limit,
		Offset:       offset,
		ServiceTypes: serviceTypes,
	})
	if err != nil {
		var status = 500
		switch {
		case errors.Is(err, model.ErrTenderNotFound):
			status = 404
		case errors.Is(err, model.ErrInternal):
			status = 500
		}
		helper.Respond(ctx, w, status, dto.NewErrResponse(err))
		return
	}
	helper.Respond(r.Context(), w, 200, result)
}

func (h *Handler) CreateTender(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tender, err := helper.ParseTenderFromBody(r)
	if err != nil {
		helper.Respond(ctx, w, 400, dto.NewErrResponse(err))
		return
	}
	createdTender, err := h.uc.CreateTender(ctx, tender)
	if err != nil {
		var status = 500
		switch {
		case errors.Is(err, model.ErrUserNotFound):
			status = 401
		case errors.Is(err, model.ErrNoRights):
			status = 403
		case errors.Is(err, model.ErrInternal):
			status = 500
		}
		helper.Respond(ctx, w, status, dto.NewErrResponse(err))
		return
	}
	helper.Respond(r.Context(), w, 200, createdTender)
}

func (h *Handler) GetMyTenders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limit, offset, err := helper.ParseLimitOffset(r)
	if err != nil {
		helper.Respond(ctx, w, 400, dto.NewErrResponse(err))
		return
	}
	username := helper.ParseUsername(r)
	tenders, err := h.uc.GetMyTenders(ctx, usecase.GetMyTendersInput{
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

func (h *Handler) GetTenderStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenderID := helper.ParseTenderID(r)
	username := helper.ParseUsername(r)
	stat, err := h.uc.GetTenderStatus(ctx, usecase.GetTenderStatusInput{
		TenderID: tenderID,
		Username: username,
	})
	if err != nil {
		var status = 500
		switch {
		case errors.Is(err, model.ErrUserNotFound):
			status = 401
		case errors.Is(err, model.ErrNoRights):
			status = 403
		case errors.Is(err, model.ErrTenderNotFound):
			status = 404
		}
		helper.Respond(ctx, w, status, dto.NewErrResponse(err))
		return
	}
	helper.Respond(r.Context(), w, 200, stat)
}

func (h *Handler) UpdateTenderStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenderID := helper.ParseTenderID(r)
	username := helper.ParseUsername(r)
	st := helper.ParseStatus(r)
	updTender, err := h.uc.UpdateTenderStatus(ctx, usecase.UpdateTenderStatusInput{
		TenderID: tenderID,
		Username: username,
		Status:   st,
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
		case errors.Is(err, model.ErrTenderNotFound):
			status = 404
		}
		helper.Respond(ctx, w, status, dto.NewErrResponse(err))
		return
	}
	helper.Respond(r.Context(), w, 200, updTender)
}

func (h *Handler) UpdateTender(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenderID := helper.ParseTenderID(r)
	username := helper.ParseUsername(r)
	info, err := helper.ParseUpdateTenderInfo(r)
	if err != nil {
		helper.Respond(ctx, w, 400, dto.NewErrResponse(err))
		return
	}
	updTender, err := h.uc.UpdateTender(ctx, usecase.UpdateTenderInput{
		TenderID:    tenderID,
		Username:    username,
		Name:        info.Name,
		Description: info.Description,
		ServiceType: info.ServiceType,
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
		case errors.Is(err, model.ErrTenderNotFound):
			status = 404
		}
		helper.Respond(ctx, w, status, dto.NewErrResponse(err))
		return
	}
	helper.Respond(r.Context(), w, 200, updTender)
}

func (h *Handler) RollbackTender(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenderID := helper.ParseTenderID(r)
	username := helper.ParseUsername(r)
	version, err := helper.ParseVersion(r)
	if err != nil {
		helper.Respond(ctx, w, 400, dto.NewErrResponse(err))
		return
	}
	updTender, err := h.uc.RollbackTender(ctx, usecase.RollbackTenderInput{
		TenderID: tenderID,
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
		case errors.Is(err, model.ErrTenderNotFound) || errors.Is(err, model.ErrNoSuchVersion):
			status = 404
		}
		helper.Respond(ctx, w, status, dto.NewErrResponse(err))
		return
	}
	helper.Respond(r.Context(), w, 200, updTender)
}
