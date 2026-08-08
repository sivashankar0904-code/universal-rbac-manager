package service

import (
	"urm/internal/helper"
	"urm/internal/repo"
	"urm/internal/service/organisation"
)

type Service struct {
	OrgSvc    *organisation.Service
	repoSvc   *repo.Service
	helperSvc *helper.Service
}

func NewService(
	repoSvc *repo.Service,
	helperSvc *helper.Service,
) *Service {
	orgSvc := organisation.NewOrganisationService(repoSvc, helperSvc)
	return &Service{
		OrgSvc: orgSvc,
	}
}
