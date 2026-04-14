package repository

import (
	"context"
	"time"

	"barber-app/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AppointmentRepo struct {
	pool *pgxpool.Pool
}

func NewAppointmentRepo(pool *pgxpool.Pool) *AppointmentRepo {
	return &AppointmentRepo{pool: pool}
}

func (r *AppointmentRepo) Create(ctx context.Context, clientID int, date, timeSlot string) (*models.Appointment, error) {
	var a models.Appointment
	err := r.pool.QueryRow(ctx,
		`INSERT INTO appointments (client_id, date, time, status)
		 VALUES ($1, $2, $3, 'scheduled')
		 RETURNING id, client_id, date::text, time::text, status, created_at`,
		clientID, date, timeSlot,
	).Scan(&a.ID, &a.ClientID, &a.Date, &a.Time, &a.Status, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AppointmentRepo) CountActiveByClient(ctx context.Context, clientID int) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM appointments WHERE client_id = $1 AND status = 'scheduled' AND (date > CURRENT_DATE OR (date = CURRENT_DATE AND time > CURRENT_TIME))",
		clientID,
	).Scan(&count)
	return count, err
}

func (r *AppointmentRepo) GetByClient(ctx context.Context, clientID int) ([]models.Appointment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, client_id, date::text, time::text, status, created_at
		 FROM appointments WHERE client_id = $1 AND status = 'scheduled'
		 ORDER BY date, time`, clientID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		var a models.Appointment
		if err := rows.Scan(&a.ID, &a.ClientID, &a.Date, &a.Time, &a.Status, &a.CreatedAt); err != nil {
			return nil, err
		}
		appointments = append(appointments, a)
	}
	return appointments, nil
}

func (r *AppointmentRepo) Cancel(ctx context.Context, id, clientID int) error {
	_, err := r.pool.Exec(ctx,
		"DELETE FROM appointments WHERE id = $1 AND client_id = $2 AND status = 'scheduled'",
		id, clientID,
	)
	return err
}

// GetByID retorna um agendamento pelo ID.
func (r *AppointmentRepo) GetByID(ctx context.Context, id int) (*models.Appointment, error) {
	var a models.Appointment
	err := r.pool.QueryRow(ctx,
		"SELECT id, client_id, date::text, time::text, status, created_at FROM appointments WHERE id = $1",
		id,
	).Scan(&a.ID, &a.ClientID, &a.Date, &a.Time, &a.Status, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AppointmentRepo) GetByDate(ctx context.Context, date string) ([]models.Appointment, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT a.id, a.client_id, c.name, a.date::text, a.time::text, a.status, a.created_at
		 FROM appointments a JOIN clients c ON a.client_id = c.id
		 WHERE a.date = $1 AND a.status = 'scheduled'
		 ORDER BY a.time`, date,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		var a models.Appointment
		if err := rows.Scan(&a.ID, &a.ClientID, &a.ClientName, &a.Date, &a.Time, &a.Status, &a.CreatedAt); err != nil {
			return nil, err
		}
		appointments = append(appointments, a)
	}
	return appointments, nil
}

func (r *AppointmentRepo) AdminCancel(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx,
		"DELETE FROM appointments WHERE id = $1 AND status = 'scheduled'", id,
	)
	return err
}

func (r *AppointmentRepo) GetUpcoming(ctx context.Context, reminderMinutes int) ([]models.Appointment, error) {
	target := time.Now().Add(time.Duration(reminderMinutes) * time.Minute)
	rows, err := r.pool.Query(ctx,
		`SELECT a.id, a.client_id, c.name, a.date::text, a.time::text, a.status, a.created_at
		 FROM appointments a JOIN clients c ON a.client_id = c.id
		 WHERE a.status = 'scheduled'
		   AND a.date = $1
		   AND a.time = $2`,
		target.Format("2006-01-02"), target.Format("15:04")+":00",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		var a models.Appointment
		if err := rows.Scan(&a.ID, &a.ClientID, &a.ClientName, &a.Date, &a.Time, &a.Status, &a.CreatedAt); err != nil {
			return nil, err
		}
		appointments = append(appointments, a)
	}
	return appointments, nil
}

func (r *AppointmentRepo) IsSlotTaken(ctx context.Context, date, timeSlot string) (bool, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		"SELECT COUNT(*) FROM appointments WHERE date = $1 AND time = $2 AND status = 'scheduled'",
		date, timeSlot,
	).Scan(&count)
	return count > 0, err
}
