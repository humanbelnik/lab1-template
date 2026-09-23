package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"lab1/internal/model"
)

var ErrNotFound = errors.New("person not found")

const schema = `CREATE TABLE IF NOT EXISTS persons (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	age INTEGER,
	address VARCHAR(255),
	work VARCHAR(255)
)`

type PersonRepository interface {
	Create(ctx context.Context, person model.Person) (int64, error)
	GetByID(ctx context.Context, id int64) (model.Person, error)
	GetAll(ctx context.Context) ([]model.Person, error)
	Update(ctx context.Context, person model.Person) error
	Delete(ctx context.Context, id int64) error
}

type PostgresPersonRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresPersonRepository(pool *pgxpool.Pool) *PostgresPersonRepository {
	return &PostgresPersonRepository{pool: pool}
}

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, schema)
	return err
}

func (r *PostgresPersonRepository) Create(ctx context.Context, person model.Person) (int64, error) {
	var id int64
	err := r.pool.QueryRow(ctx,
		"INSERT INTO persons (name, age, address, work) VALUES ($1, $2, $3, $4) RETURNING id",
		person.Name, person.Age, person.Address, person.Work,
	).Scan(&id)
	return id, err
}

func (r *PostgresPersonRepository) GetByID(ctx context.Context, id int64) (model.Person, error) {
	var p model.Person
	err := r.pool.QueryRow(ctx,
		"SELECT id, name, age, address, work FROM persons WHERE id = $1", id,
	).Scan(&p.ID, &p.Name, &p.Age, &p.Address, &p.Work)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Person{}, ErrNotFound
	}
	return p, err
}

func (r *PostgresPersonRepository) GetAll(ctx context.Context) ([]model.Person, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, name, age, address, work FROM persons ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	persons := make([]model.Person, 0)
	for rows.Next() {
		var p model.Person
		if err := rows.Scan(&p.ID, &p.Name, &p.Age, &p.Address, &p.Work); err != nil {
			return nil, err
		}
		persons = append(persons, p)
	}
	return persons, rows.Err()
}

func (r *PostgresPersonRepository) Update(ctx context.Context, person model.Person) error {
	tag, err := r.pool.Exec(ctx,
		"UPDATE persons SET name = $1, age = $2, address = $3, work = $4 WHERE id = $5",
		person.Name, person.Age, person.Address, person.Work, person.ID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PostgresPersonRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.pool.Exec(ctx, "DELETE FROM persons WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
