package service

import (
	"context"
	"errors"

	"lab1/internal/model"
	"lab1/internal/repository"
)

var (
	ErrNotFound   = repository.ErrNotFound
	ErrValidation = errors.New("validation error")
)

type PersonService struct {
	repo repository.PersonRepository
}

func NewPersonService(repo repository.PersonRepository) *PersonService {
	return &PersonService{repo: repo}
}

func (s *PersonService) Create(ctx context.Context, req model.CreatePersonRequest) (int64, error) {
	if req.Name == "" {
		return 0, ErrValidation
	}
	person := model.Person{
		Name:    req.Name,
		Age:     req.Age,
		Address: req.Address,
		Work:    req.Work,
	}
	return s.repo.Create(ctx, person)
}

func (s *PersonService) GetByID(ctx context.Context, id int64) (model.Person, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PersonService) GetAll(ctx context.Context) ([]model.Person, error) {
	return s.repo.GetAll(ctx)
}

func (s *PersonService) Update(ctx context.Context, id int64, req model.UpdatePersonRequest) (model.Person, error) {
	person, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return model.Person{}, err
	}

	if req.Name != nil {
		if *req.Name == "" {
			return model.Person{}, ErrValidation
		}
		person.Name = *req.Name
	}
	if req.Age != nil {
		person.Age = *req.Age
	}
	if req.Address != nil {
		person.Address = *req.Address
	}
	if req.Work != nil {
		person.Work = *req.Work
	}

	if err := s.repo.Update(ctx, person); err != nil {
		return model.Person{}, err
	}
	return person, nil
}

func (s *PersonService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
