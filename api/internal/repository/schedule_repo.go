package repository

import (
	"context"

	"barber-app/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ScheduleRepo struct {
	pool *pgxpool.Pool
}

func NewScheduleRepo(pool *pgxpool.Pool) *ScheduleRepo {
	return &ScheduleRepo{pool: pool}
}

func (r *ScheduleRepo) Upsert(ctx context.Context, s models.Schedule) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO schedules (day_of_week, slot, start_time, end_time, active)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (day_of_week, slot) DO UPDATE SET start_time = $3, end_time = $4, active = $5`,
		s.DayOfWeek, s.Slot, s.StartTime, s.EndTime, s.Active,
	)
	return err
}

func (r *ScheduleRepo) GetAll(ctx context.Context) ([]models.Schedule, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, day_of_week, start_time::text, end_time::text, active, slot FROM schedules ORDER BY day_of_week, slot",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var s models.Schedule
		if err := rows.Scan(&s.ID, &s.DayOfWeek, &s.StartTime, &s.EndTime, &s.Active, &s.Slot); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func (r *ScheduleRepo) GetByDay(ctx context.Context, dayOfWeek int) ([]models.Schedule, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, day_of_week, start_time::text, end_time::text, active, slot FROM schedules WHERE day_of_week = $1 AND active = true ORDER BY slot",
		dayOfWeek,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var s models.Schedule
		if err := rows.Scan(&s.ID, &s.DayOfWeek, &s.StartTime, &s.EndTime, &s.Active, &s.Slot); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}
