package repo

import (
	"urm/internal/helper"
	organisationrepo "urm/internal/repo/organisation_repo"
)

type Service struct {
	db        helper.Querier
	OrgRepo   *organisationrepo.Service
	helperSvc *helper.Service
}

func NewService(
	db helper.Querier,
	helperSvc *helper.Service,
) *Service {
	return &Service{
		db:        db,
		OrgRepo:   organisationrepo.NewService(db, helperSvc),
		helperSvc: helperSvc,
	}
}
