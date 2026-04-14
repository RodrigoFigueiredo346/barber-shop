package repository

import (
	"context"

	"barber-app/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BlockedSlotRepo struct {
	pool *pgxpool.Pool
}

func NewBlockedSlotRepo(pool *pgxpool.Pool) *BlockedSlotRepo {
	return &BlockedSlotRepo{pool: pool}
}

func (r *BlockedSlotRepo) Block(ctx context.Context, date, timeSlot string) error {
	_, err := r.pool.Exec(ctx,
		"INSERT INTO blocked_slots (date, time) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		date, timeSlot,
	)
	return err
}

func (r *BlockedSlotRepo) Unblock(ctx context.Context, date, timeSlot string) error {
	_, err := r.pool.Exec(ctx,
		"DELETE FROM blocked_slots WHERE date = $1 AND time = $2",
		date, timeSlot,
	)
	return err
}

func (r *BlockedSlotRepo) GetByDate(ctx context.Context, date string) ([]models.BlockedSlot, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, date::text, time::text FROM blocked_slots WHERE date = $1", date,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []models.BlockedSlot
	for rows.Next() {
		var s models.BlockedSlot
		if err := rows.Scan(&s.ID, &s.Date, &s.Time); err != nil {
			return nil, err
		}
		slots = append(slots, s)
	}
	return slots, nil
}

func (r *BlockedSlotRepo) IsBlocked(ctx context.Context, date, timeSlot string) (bool, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM blocked_slots WHERE date = $1 AND time = $2",
		date, timeSlot,
	).Scan(&count)
	return count > 0, err
}
