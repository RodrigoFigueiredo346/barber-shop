package repository

import (
	"context"

	"barber-app/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SettingsRepo struct {
	pool *pgxpool.Pool
}

func NewSettingsRepo(pool *pgxpool.Pool) *SettingsRepo {
	return &SettingsRepo{pool: pool}
}

func (r *SettingsRepo) Get(ctx context.Context) (*models.Settings, error) {
	var s models.Settings
	err := r.pool.QueryRow(ctx,
		"SELECT reminder_minutes FROM settings WHERE id = 1",
	).Scan(&s.ReminderMinutes)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SettingsRepo) Update(ctx context.Context, s models.Settings) error {
	_, err := r.pool.Exec(ctx,
		"UPDATE settings SET reminder_minutes = $1 WHERE id = 1",
		s.ReminderMinutes,
	)
	return err
}
