package organisationrepo

import (
	"context"

	"urm/internal/apperr"
	"urm/internal/helper"
	"urm/internal/model"
)

type Service struct {
	db        helper.Querier
	helperSvc *helper.Service
}

func NewService(
	db helper.Querier,
	helperSvc *helper.Service,
) *Service {
	return &Service{
		db:        db,
		helperSvc: helperSvc,
	}
}

type OrganisationRepo interface {
	Create(ctx context.Context, o *model.Organisation) error
	Get(ctx context.Context, id string) (*model.Organisation, error)
	List(ctx context.Context) ([]*model.Organisation, error)
	Update(ctx context.Context, o *model.Organisation) error
	SetActive(ctx context.Context, id string, active bool) (*model.Organisation, error)
	Delete(ctx context.Context, id string) error
}

func (r *Service) Create(ctx context.Context, o *model.Organisation) error {

	const query = `
		INSERT INTO organisation (key, name, description, parent_id, is_active, meta)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		o.Key, o.Name, o.Description, o.ParentID, o.IsActive, o.Meta,
	).Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return r.helperSvc.TranslateErr(err)
	}
	return nil
}

func (r *Service) Get(ctx context.Context, id string) (*model.Organisation, error) {
	const query = `
		SELECT id, key, name, description, parent_id, is_active, meta, created_at, updated_at
		FROM organisation
		WHERE id = $1`

	o := &model.Organisation{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&o.ID, &o.Key, &o.Name, &o.Description, &o.ParentID, &o.IsActive, &o.Meta, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, r.helperSvc.TranslateErr(err)
	}
	return o, nil
}

func (r *Service) List(ctx context.Context) ([]*model.Organisation, error) {
	const query = `
		SELECT id, key, name, description, parent_id, is_active, meta, created_at, updated_at
		FROM organisation
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, r.helperSvc.TranslateErr(err)
	}
	defer rows.Close()

	orgs := make([]*model.Organisation, 0)
	for rows.Next() {
		o := &model.Organisation{}
		if err := rows.Scan(
			&o.ID, &o.Key, &o.Name, &o.Description, &o.ParentID, &o.IsActive, &o.Meta, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, r.helperSvc.TranslateErr(err)
		}
		orgs = append(orgs, o)
	}
	if err := rows.Err(); err != nil {
		return nil, r.helperSvc.TranslateErr(err)
	}
	return orgs, nil
}

// Update persists the generically-editable fields only (name, description).
// IsActive goes through SetActive instead.
func (r *Service) Update(ctx context.Context, o *model.Organisation) error {
	const query = `
		UPDATE organisation
		SET name = $2, description = $3, updated_at = now()
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRow(ctx, query, o.ID, o.Name, o.Description).Scan(&o.UpdatedAt)
	if err != nil {
		return r.helperSvc.TranslateErr(err)
	}
	return nil
}

// SetActive flips is_active and returns the full updated row — used by
// service.Onboard (and, later, any deactivate/suspend flow).
func (r *Service) SetActive(ctx context.Context, id string, active bool) (*model.Organisation, error) {
	const query = `
		UPDATE organisation
		SET is_active = $2, updated_at = now()
		WHERE id = $1
		RETURNING id, key, name, description, parent_id, is_active, meta, created_at, updated_at`

	o := &model.Organisation{}
	err := r.db.QueryRow(ctx, query, id, active).Scan(
		&o.ID, &o.Key, &o.Name, &o.Description, &o.ParentID, &o.IsActive, &o.Meta, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, r.helperSvc.TranslateErr(err)
	}
	return o, nil
}

func (r *Service) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM organisation WHERE id = $1`

	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return r.helperSvc.TranslateErr(err)
	}
	if tag.RowsAffected() == 0 {
		return apperr.NotFound("organisation not found")
	}
	return nil
}
