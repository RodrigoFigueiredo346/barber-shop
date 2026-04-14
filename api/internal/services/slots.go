package services

import (
	"context"
	"fmt"
	"log"
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

// GetAvailableSlots retorna os horários disponíveis para uma data.
func (s *SlotService) GetAvailableSlots(ctx context.Context, date string) ([]string, error) {
	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("data inválida: %w", err)
	}

	dayOfWeek := int(d.Weekday())
	log.Printf("[SLOTS] date=%s dayOfWeek=%d", date, dayOfWeek)
	schedule, err := s.scheduleRepo.GetByDay(ctx, dayOfWeek)
	if err != nil {
		log.Printf("[SLOTS] Sem schedule para dayOfWeek=%d: %v", dayOfWeek, err)
		return nil, nil // dia sem configuração = sem horários
	}
	log.Printf("[SLOTS] Schedule encontrado: active=%v start=%s end=%s", schedule.Active, schedule.StartTime, schedule.EndTime)
	if !schedule.Active {
		log.Printf("[SLOTS] Dia %d não está ativo", dayOfWeek)
		return nil, nil
	}

	startStr := schedule.StartTime
	endStr := schedule.EndTime
	log.Printf("[SLOTS] Parsing start=%q end=%q", startStr, endStr)
	// PostgreSQL retorna TIME como HH:MM:SS, precisamos aceitar ambos formatos
	start, err := time.Parse("15:04:05", startStr)
	if err != nil {
		start, err = time.Parse("15:04", startStr)
		if err != nil {
			return nil, nil
		}
	}
	end, err2 := time.Parse("15:04:05", endStr)
	if err2 != nil {
		end, err2 = time.Parse("15:04", endStr)
		if err2 != nil {
			return nil, nil
		}
	}

	var allSlots []string
	for t := start; t.Before(end); t = t.Add(30 * time.Minute) {
		allSlots = append(allSlots, t.Format("15:04"))
	}

	log.Printf("[SLOTS] Total slots gerados: %d", len(allSlots))
	log.Printf("[SLOTS] allSlots: %v", allSlots)

	blocked, _ := s.blockedSlotRepo.GetByDate(ctx, date)
	blockedMap := make(map[string]bool)
	for _, b := range blocked {
		// Normaliza pra HH:MM (banco pode retornar HH:MM:SS)
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
		taken, _ := s.appointmentRepo.IsSlotTaken(ctx, date, slot+":00")
		if !taken {
			available = append(available, slot)
		}
	}

	return available, nil
}
