package service_test

import (
	"context"
	"errors"
	"testing"

	"lab1/internal/model"
	"lab1/internal/repository"
	"lab1/internal/service"
)

type fakeRepository struct {
	persons map[int64]model.Person
	nextID  int64
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{persons: make(map[int64]model.Person)}
}

func (r *fakeRepository) Create(_ context.Context, person model.Person) (int64, error) {
	r.nextID++
	person.ID = r.nextID
	r.persons[person.ID] = person
	return person.ID, nil
}

func (r *fakeRepository) GetByID(_ context.Context, id int64) (model.Person, error) {
	p, ok := r.persons[id]
	if !ok {
		return model.Person{}, repository.ErrNotFound
	}
	return p, nil
}

func (r *fakeRepository) GetAll(_ context.Context) ([]model.Person, error) {
	persons := make([]model.Person, 0, len(r.persons))
	for _, p := range r.persons {
		persons = append(persons, p)
	}
	return persons, nil
}

func (r *fakeRepository) Update(_ context.Context, person model.Person) error {
	if _, ok := r.persons[person.ID]; !ok {
		return repository.ErrNotFound
	}
	r.persons[person.ID] = person
	return nil
}

func (r *fakeRepository) Delete(_ context.Context, id int64) error {
	if _, ok := r.persons[id]; !ok {
		return repository.ErrNotFound
	}
	delete(r.persons, id)
	return nil
}

func TestPersonService_Create(t *testing.T) {
	svc := service.NewPersonService(newFakeRepository())

	id, err := svc.Create(context.Background(), model.CreatePersonRequest{Name: "Alice", Age: 30})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == 0 {
		t.Fatal("expected non-zero id")
	}
}

func TestPersonService_Create_ValidationError(t *testing.T) {
	svc := service.NewPersonService(newFakeRepository())

	_, err := svc.Create(context.Background(), model.CreatePersonRequest{Name: ""})
	if !errors.Is(err, service.ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestPersonService_GetByID_NotFound(t *testing.T) {
	svc := service.NewPersonService(newFakeRepository())

	_, err := svc.GetByID(context.Background(), 42)
	if !errors.Is(err, service.ErrNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestPersonService_Update(t *testing.T) {
	svc := service.NewPersonService(newFakeRepository())

	id, err := svc.Create(context.Background(), model.CreatePersonRequest{Name: "Bob", Work: "Acme"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	newName := "Bobby"
	updated, err := svc.Update(context.Background(), id, model.UpdatePersonRequest{Name: &newName})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Name != "Bobby" || updated.Work != "Acme" {
		t.Fatalf("unexpected person after update: %+v", updated)
	}
}

func TestPersonService_Delete(t *testing.T) {
	svc := service.NewPersonService(newFakeRepository())

	id, err := svc.Create(context.Background(), model.CreatePersonRequest{Name: "Carol"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := svc.Delete(context.Background(), id); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := svc.GetByID(context.Background(), id); !errors.Is(err, service.ErrNotFound) {
		t.Fatal("expected not found after delete")
	}
}
