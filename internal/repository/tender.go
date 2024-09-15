package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/b0pof/avito-internship/internal/model"
	"github.com/b0pof/avito-internship/pkg/logger"
)

func (r *Repository) GetTenderByID(ctx context.Context, tenderID string) (model.Tender, error) {
	q := `SELECT t.id, tv.name, tv.description, t.status, tv.service_type, tv.version, t.created_at
		FROM tender_version tv
			INNER JOIN tender t on tv.tender_id = t.id
		WHERE tender_id = $1
			AND version = (SELECT MAX(version) FROM tender_version WHERE tender_id = $1);`

	var tender model.Tender
	if err := r.db.GetContext(ctx, &tender, q, tenderID); err != nil {
		return model.Tender{}, model.ErrTenderNotFound
	}
	return tender, nil
}

type GetTendersInput struct {
	ServiceTypes []string
	Limit        int
	Offset       int
}

func (r *Repository) GetTenders(ctx context.Context, input GetTendersInput) ([]model.Tender, error) {
	var q string
	tenders := make([]model.Tender, 0)

	if len(input.ServiceTypes) == 0 {
		q = `SELECT t.id, tv.name, tv.description, t.status, tv.service_type, tv.version, t.created_at
				FROM tender t
				INNER JOIN tender_version tv
					ON t.id = tv.tender_id
					AND tv.version = (
						SELECT MAX(version)
						FROM tender_version
						WHERE tender_id = t.id
					)
				WHERE status = 'Published'
				ORDER BY tv.name
				LIMIT $1
				OFFSET $2;`
		if err := r.db.SelectContext(ctx, &tenders, q, input.Limit, input.Offset); err != nil {
			logger.Error(ctx, err.Error())
			return nil, model.ErrInternal
		}
	} else {
		q = `SELECT t.id, tv.name, tv.description, t.status, tv.service_type, tv.version, t.created_at
				FROM tender t
				INNER JOIN tender_version tv
					ON t.id = tv.tender_id
					AND tv.version = (
						SELECT MAX(version)
						FROM tender_version
						WHERE tender_id = t.id
					)
				WHERE service_type = ANY($1::tender_service_type[]) AND status = 'Published'
				ORDER BY tv.name
				LIMIT $2
				OFFSET $3;`
		if err := r.db.SelectContext(ctx, &tenders, q, input.ServiceTypes, input.Limit, input.Offset); err != nil {
			logger.Error(ctx, err.Error())
			return nil, model.ErrInternal
		}
	}
	return tenders, nil
}

type CreateTenderInput struct {
	Name           string
	Description    string
	ServiceType    string
	OrganizationID string
	CreatorID      string
}

func (r *Repository) CreateTender(ctx context.Context, input CreateTenderInput) (model.Tender, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return model.Tender{}, model.ErrInternal
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	q := `INSERT INTO tender (organization_id, author_id)
			VALUES ($1, $2)
			RETURNING id;`

	var tenderID string
	if err = tx.GetContext(ctx, &tenderID, q, input.OrganizationID, input.CreatorID); err != nil {
		logger.Error(ctx, err.Error())
		return model.Tender{}, model.ErrInternal
	}

	q = `INSERT INTO tender_version (tender_id, name, description, service_type, version)
			SELECT $1, $2, $3, $4, COALESCE(tv.version, 0) + 1
			FROM tender t
				LEFT JOIN tender_version tv ON t.id = tv.tender_id
			WHERE t.id = $1
			ORDER BY version DESC
			LIMIT 1;`

	if _, err = tx.ExecContext(ctx, q, tenderID, input.Name, input.Description, input.ServiceType); err != nil {
		logger.Error(ctx, err.Error())
		return model.Tender{}, model.ErrInternal
	}
	_ = tx.Commit()
	return r.GetTenderByID(ctx, tenderID)
}

type EditTenderInput struct {
	TenderID    string
	Name        string
	Description string
	ServiceType string
}

func (r *Repository) UpdateTender(ctx context.Context, input EditTenderInput) (model.Tender, error) {
	q := `INSERT INTO tender_version (tender_id, name, description, service_type, version)
			SELECT t.id, %s, %s, %s, COALESCE(tv.version, 0) + 1
			FROM tender t
				LEFT JOIN tender_version tv ON t.id = tv.tender_id
			WHERE t.id = $1
			ORDER BY version DESC
			LIMIT 1;`

	format := "'%s'"

	var name = "name"
	if input.Name != "" {
		name = fmt.Sprintf(format, input.Name)
	}
	var description = "description"
	if input.Description != "" {
		description = fmt.Sprintf(format, input.Description)
	}
	var serviceType = "service_type"
	if input.ServiceType != "" {
		serviceType = fmt.Sprintf(format, input.ServiceType)
	}
	q = fmt.Sprintf(q, name, description, serviceType)

	if _, err := r.db.ExecContext(ctx, q, input.TenderID); err != nil {
		if strings.Contains(err.Error(), "invalid value") {
			return model.Tender{}, model.ErrInvalidAttributeValue
		}
		logger.Error(ctx, err.Error())
		return model.Tender{}, model.ErrInternal
	}
	return r.GetTenderByID(ctx, input.TenderID)
}

type TenderHasVersionInput struct {
	TenderID string
	Version  int
}

func (r *Repository) TenderHasVersion(ctx context.Context, input TenderHasVersionInput) (bool, error) {
	q := `SELECT (SELECT tender_id
			FROM tender_version
			WHERE tender_id = $1
			  AND version = $2
		) IS NOT NULL;`

	var hasVersion bool
	if err := r.db.GetContext(ctx, &hasVersion, q, input.TenderID, input.Version); err != nil {
		logger.Error(ctx, err.Error())
		return false, model.ErrInternal
	}
	return hasVersion, nil
}

func (r *Repository) TenderExists(ctx context.Context, tenderID string) bool {
	q := `SELECT id FROM tender WHERE id = $1;`

	var foundTenderID string
	if err := r.db.GetContext(ctx, &foundTenderID, q, tenderID); err != nil {
		logger.Error(ctx, err.Error())
		return false
	}
	return foundTenderID == tenderID
}

type RollbackTenderInput struct {
	TenderID string
	Version  int
}

func (r *Repository) RollbackTender(ctx context.Context, input RollbackTenderInput) (model.Tender, error) {
	q := `INSERT INTO tender_version (tender_id, name, description, service_type, version)
			SELECT t.id, name, description, service_type, COALESCE(
				(SELECT MAX(version) FROM tender_version WHERE tender_id = $1), 0
			) + 1
			FROM tender t
				LEFT JOIN tender_version tv ON t.id = tv.tender_id
			WHERE t.id = $1 AND tv.version = $2;`

	if _, err := r.db.ExecContext(ctx, q, input.TenderID, input.Version); err != nil {
		logger.Error(ctx, err.Error())
		return model.Tender{}, model.ErrInternal
	}
	return r.GetTenderByID(ctx, input.TenderID)
}

type GetMyTendersInput struct {
	Limit  int
	Offset int
	UserID string
}

func (r *Repository) GetMyTenders(ctx context.Context, input GetMyTendersInput) ([]model.Tender, error) {
	q := `SELECT t.id, tv.name, tv.description, t.status, tv.service_type, tv.version, t.created_at
			FROM tender t
			INNER JOIN tender_version tv
				ON t.id = tv.tender_id
				AND tv.version = (
					SELECT MAX(version)
					FROM tender_version
					WHERE tender_id = t.id
				)
			WHERE t.author_id = $1
			ORDER BY tv.name
			LIMIT $2
			OFFSET $3;`

	tenders := make([]model.Tender, 0)
	if err := r.db.SelectContext(ctx, &tenders, q, input.UserID, input.Limit, input.Offset); err != nil {
		logger.Error(ctx, err.Error())
		return nil, model.ErrInternal
	}
	return tenders, nil
}

func (r *Repository) GetTenderStatus(ctx context.Context, tenderID string) (string, error) {
	q := `SELECT status FROM tender WHERE id = $1`

	var status string
	if err := r.db.GetContext(ctx, &status, q, tenderID); err != nil {
		return "", model.ErrTenderNotFound
	}
	return status, nil
}

func (r *Repository) CloseTenderByBidID(ctx context.Context, bidID string) error {
	q := `UPDATE tender
		SET status = 'Closed'
		WHERE id = (
			SELECT tender_id
			FROM bid
			WHERE id = $1
		);`

	if _, err := r.db.ExecContext(ctx, q, bidID); err != nil {
		return model.ErrTenderNotFound
	}
	return nil
}

type UpdateTenderStatusInput struct {
	TenderID string
	Status   string
}

func (r *Repository) UpdateTenderStatus(ctx context.Context, input UpdateTenderStatusInput) (model.Tender, error) {
	q := `UPDATE tender
			SET status = $1
			WHERE id = $2;`

	if _, err := r.db.ExecContext(ctx, q, input.Status, input.TenderID); err != nil {
		if strings.Contains(err.Error(), "invalid input") {
			return model.Tender{}, model.ErrInvalidAttributeValue
		}
		return model.Tender{}, model.ErrTenderNotFound
	}
	return r.GetTenderByID(ctx, input.TenderID)
}

func (r *Repository) IsTenderExist(ctx context.Context, tenderID string) bool {
	q := `SELECT id FROM tender WHERE id = $1`

	var foundTenderID string
	if err := r.db.GetContext(ctx, &foundTenderID, q, tenderID); err != nil {
		logger.Error(ctx, err.Error())
		return false
	}
	return foundTenderID == tenderID
}
