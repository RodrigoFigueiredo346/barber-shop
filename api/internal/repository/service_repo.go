package repository

import (
	"context"

	"barber-app/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRepo struct {
	pool *pgxpool.Pool
}

func NewServiceRepo(pool *pgxpool.Pool) *ServiceRepo {
	return &ServiceRepo{pool: pool}
}

func (r *ServiceRepo) GetAll(ctx context.Context) ([]models.Service, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, name, duration, price, active FROM services ORDER BY name",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []models.Service
	for rows.Next() {
		var s models.Service
		if err := rows.Scan(&s.ID, &s.Name, &s.Duration, &s.Price, &s.Active); err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}

func (r *ServiceRepo) GetActive(ctx context.Context) ([]models.Service, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, name, duration, price, active FROM services WHERE active = true ORDER BY name",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []models.Service
	for rows.Next() {
		var s models.Service
		if err := rows.Scan(&s.ID, &s.Name, &s.Duration, &s.Price, &s.Active); err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}

func (r *ServiceRepo) Create(ctx context.Context, s models.Service) (*models.Service, error) {
	var created models.Service
	err := r.pool.QueryRow(ctx,
		"INSERT INTO services (name, duration, price, active) VALUES ($1, $2, $3, $4) RETURNING id, name, duration, price, active",
		s.Name, s.Duration, s.Price, s.Active,
	).Scan(&created.ID, &created.Name, &created.Duration, &created.Price, &created.Active)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (r *ServiceRepo) Update(ctx context.Context, s models.Service) error {
	_, err := r.pool.Exec(ctx,
		"UPDATE services SET name = $1, duration = $2, price = $3, active = $4 WHERE id = $5",
		s.Name, s.Duration, s.Price, s.Active, s.ID,
	)
	return err
}

func (r *ServiceRepo) Delete(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM services WHERE id = $1", id)
	return err
}
