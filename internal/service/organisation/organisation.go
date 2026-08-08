package organisation

import (
	"context"
	"urm/internal/apperr"
	"urm/internal/helper"
	"urm/internal/model"
	"urm/internal/repo"
)

type Service struct {
	repoSvc   *repo.Service
	helperSvc *helper.Service
}

func NewOrganisationService(
	repoSvc *repo.Service,
	helperSvc *helper.Service,
) *Service {
	return &Service{
		repoSvc:   repoSvc,
		helperSvc: helperSvc,
	}
}

func (s *Service) Create(ctx context.Context, o model.Organisation) (*model.Organisation, error) {

	res, err := s.helperSvc.StructToMapStringWithTag(o)
	if err != nil {
		return nil, err
	}
	validation_fields := []string{"key", "name"}

	if err := s.helperSvc.ValidateRequiredFields(res, validation_fields); err != nil {
		return nil, err
	}

	o.IsActive = false // starts pending — Onboard activates it
	if err := s.repoSvc.OrgRepo.Create(ctx, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

func (s *Service) Get(ctx context.Context, id string) (*model.Organisation, error) {
	return s.repoSvc.OrgRepo.Get(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]*model.Organisation, error) {
	return s.repoSvc.OrgRepo.List(ctx)
}

func (s *Service) Update(ctx context.Context, id string, p model.Organisation) (*model.Organisation, error) {
	res, err := s.helperSvc.StructToMapStringWithTag(p)
	if err != nil {
		return nil, err
	}
	validation_fields := []string{"name"}

	if err := s.helperSvc.ValidateRequiredFields(res, validation_fields); err != nil {
		return nil, err
	}

	o, err := s.repoSvc.OrgRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	o.Name = p.Name
	o.Description = p.Description
	if err := s.repoSvc.OrgRepo.Update(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repoSvc.OrgRepo.Delete(ctx, id)
}

func (s *Service) Onboard(ctx context.Context, id string) (*model.Organisation, error) {
	o, err := s.repoSvc.OrgRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if o.IsActive {
		return nil, apperr.Conflict("organisation is already active")
	}

	return s.repoSvc.OrgRepo.SetActive(ctx, id, true)
}
