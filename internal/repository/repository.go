package repository

import (
	"context"

	"github.com/b0pof/avito-internship/internal/model"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

type IRepository interface {
	ITenderRepository
	IBidRepository
	IOrganizationRepository
	IUserRepository
}

type ITenderRepository interface {
	GetTenderByID(ctx context.Context, tenderID string) (model.Tender, error)
	GetTenders(ctx context.Context, input GetTendersInput) ([]model.Tender, error)
	CreateTender(ctx context.Context, input CreateTenderInput) (model.Tender, error)
	UpdateTender(ctx context.Context, input EditTenderInput) (model.Tender, error)
	TenderHasVersion(ctx context.Context, input TenderHasVersionInput) (bool, error)
	TenderExists(ctx context.Context, tenderID string) bool
	RollbackTender(ctx context.Context, input RollbackTenderInput) (model.Tender, error)
	GetMyTenders(ctx context.Context, input GetMyTendersInput) ([]model.Tender, error)
	GetTenderStatus(ctx context.Context, tenderID string) (string, error)
	CloseTenderByBidID(ctx context.Context, bidID string) error
	UpdateTenderStatus(ctx context.Context, input UpdateTenderStatusInput) (model.Tender, error)
	IsTenderExist(ctx context.Context, tenderID string) bool
}

type IBidRepository interface {
	GetBidByID(ctx context.Context, bidID string) (model.Bid, error)
	CreateBid(ctx context.Context, input CreateBidInput) (model.Bid, error)
	GetMyBids(ctx context.Context, input GetMyBidsInput) ([]model.Bid, error)
	BidExists(ctx context.Context, bidID string) bool
	GetTenderBids(ctx context.Context, input GetTenderBidsInput) ([]model.Bid, error)
	GetBidStatus(ctx context.Context, bidID string) (string, error)
	UpdateBidStatus(ctx context.Context, input UpdateBidStatusInput) (model.Bid, error)
	UpdateBid(ctx context.Context, input EditBidInput) (model.Bid, error)
	BidHasVersion(ctx context.Context, input BidHasVersionInput) (bool, error)
	RollbackBid(ctx context.Context, input RollbackBidInput) (model.Bid, error)
}

type IOrganizationRepository interface {
	GetOrganizationIDByEmployeeID(ctx context.Context, employeeID string) (string, error)
}

type IUserRepository interface {
	IsUserOrganizationResponsible(ctx context.Context, userID, orgID string) bool
	GetUserIDByUsername(ctx context.Context, username string) (string, error)
	IsUserResponsibleForTender(ctx context.Context, tenderID, userID string) bool
	IsBidVisibleForUser(ctx context.Context, userID, bidID string) (bool, error)
	UserCanSubmitDecision(ctx context.Context, bidID, userID string) bool
	GetUserIDByBidID(ctx context.Context, bidID string) (string, error)
	IsUserExist(ctx context.Context, userID string) bool
}
