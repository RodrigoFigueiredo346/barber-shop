package services

import (
	"context"
	"fmt"
	"time"

	"barber-app/internal/repository"
)

type SlotService struct {
	scheduleRepo    *repository.ScheduleRepo
	blockedSlotRepo *repository.BlockedSlotRepo
	appointmentRepo *repository.AppointmentRepo
}

func NewSlotService(sr *repository.ScheduleRepo, br *repository.BlockedSlotRepo, ar *repository.AppointmentRepo) *SlotService {
	return &SlotService{scheduleRepo: sr, blockedSlotRepo: br, appointmentRepo: ar}
}

func parseTime(s string) (time.Time, error) {
	t, err := time.Parse("15:04:05", s)
	if err != nil {
		t, err = time.Parse("15:04", s)
	}
	return t, err
}

// GetAvailableSlots retorna os horários disponíveis para uma data.
func (s *SlotService) GetAvailableSlots(ctx context.Context, date string) ([]string, error) {
	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("data inválida: %w", err)
	}

	dayOfWeek := int(d.Weekday())
	schedules, err := s.scheduleRepo.GetByDay(ctx, dayOfWeek)
	if err != nil || len(schedules) == 0 {
		return nil, nil
	}

	// Gera slots de todas as faixas ativas
	var allSlots []string
	for _, sched := range schedules {
		start, err1 := parseTime(sched.StartTime)
		end, err2 := parseTime(sched.EndTime)
		if err1 != nil || err2 != nil {
			continue
		}
		for t := start; t.Before(end); t = t.Add(30 * time.Minute) {
			allSlots = append(allSlots, t.Format("15:04"))
		}
	}

	if len(allSlots) == 0 {
		return nil, nil
	}

	// Filtra horários passados se for o dia atual
	now := time.Now()
	isToday := d.Year() == now.Year() && d.Month() == now.Month() && d.Day() == now.Day()
	nowMinutes := now.Hour()*60 + now.Minute()

	blocked, _ := s.blockedSlotRepo.GetByDate(ctx, date)
	blockedMap := make(map[string]bool)
	for _, b := range blocked {
		t := b.Time
		if len(t) > 5 {
			t = t[:5]
		}
		blockedMap[t] = true
	}

	var available []string
	for _, slot := range allSlots {
		if blockedMap[slot] {
			continue
		}
		// Filtra passados no dia atual
		if isToday {
			parsed, err := parseTime(slot)
			if err == nil {
				slotMinutes := parsed.Hour()*60 + parsed.Minute()
				if slotMinutes <= nowMinutes {
					continue
				}
			}
		}
		taken, _ := s.appointmentRepo.IsSlotTaken(ctx, date, slot+":00")
		if !taken {
			available = append(available, slot)
		}
	}

	return available, nil
}
