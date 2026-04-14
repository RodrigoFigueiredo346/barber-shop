package repository

import (
	"context"

	"barber-app/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientRepo struct {
	pool *pgxpool.Pool
}

func NewClientRepo(pool *pgxpool.Pool) *ClientRepo {
	return &ClientRepo{pool: pool}
}

func (r *ClientRepo) FindByPhone(ctx context.Context, phone string) (*models.Client, error) {
	var c models.Client
	err := r.pool.QueryRow(ctx,
		"SELECT id, name, phone, created_at FROM clients WHERE phone = $1", phone,
	).Scan(&c.ID, &c.Name, &c.Phone, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ClientRepo) Create(ctx context.Context, name, phone string) (*models.Client, error) {
	var c models.Client
	err := r.pool.QueryRow(ctx,
		"INSERT INTO clients (name, phone) VALUES ($1, $2) RETURNING id, name, phone, created_at",
		name, phone,
	).Scan(&c.ID, &c.Name, &c.Phone, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ClientRepo) GetPhoneByID(ctx context.Context, id int) (string, error) {
	var phone string
	err := r.pool.QueryRow(ctx, "SELECT phone FROM clients WHERE id = $1", id).Scan(&phone)
	return phone, err
}
