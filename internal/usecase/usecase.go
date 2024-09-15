package usecase

import (
	"context"

	"github.com/b0pof/avito-internship/internal/model"
	"github.com/b0pof/avito-internship/internal/repository"
)

type Usecase struct {
	repo repository.IRepository
}

func New(repo repository.IRepository) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

type IUsecase interface {
	IBidUsecase
	ITenderUsecase
}

type IBidUsecase interface {
	CreateBid(ctx context.Context, input repository.CreateBidInput) (model.Bid, error)
	GetMyBids(ctx context.Context, input GetMyBidsInput) ([]model.Bid, error)
	GetTenderBids(ctx context.Context, input repository.GetTenderBidsInput) ([]model.Bid, error)
	GetBidStatus(ctx context.Context, input GetBidStatusInput) (string, error)
	UpdateBidStatus(ctx context.Context, input UpdateBidStatusInput) (model.Bid, error)
	SubmitDecision(ctx context.Context, input SubmitDecisionInput) (model.Bid, error)
	UpdateBid(ctx context.Context, input UpdateBidInput) (model.Bid, error)
	RollbackBid(ctx context.Context, input RollbackBidInput) (model.Bid, error)
}

type ITenderUsecase interface {
	GetTenders(ctx context.Context, input repository.GetTendersInput) ([]model.Tender, error)
	CreateTender(ctx context.Context, input CreateTenderInput) (model.Tender, error)
	GetMyTenders(ctx context.Context, input GetMyTendersInput) ([]model.Tender, error)
	GetTenderStatus(ctx context.Context, input GetTenderStatusInput) (string, error)
	UpdateTenderStatus(ctx context.Context, input UpdateTenderStatusInput) (model.Tender, error)
	UpdateTender(ctx context.Context, input UpdateTenderInput) (model.Tender, error)
	RollbackTender(ctx context.Context, input RollbackTenderInput) (model.Tender, error)
}
