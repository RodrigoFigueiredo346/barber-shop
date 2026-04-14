package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"barber-app/internal/repository"
	"barber-app/internal/whatsapp"
)

type ReminderService struct {
	appointmentRepo *repository.AppointmentRepo
	clientRepo      *repository.ClientRepo
	settingsRepo    *repository.SettingsRepo
	sender          whatsapp.Sender
}

func NewReminderService(ar *repository.AppointmentRepo, cr *repository.ClientRepo, sr *repository.SettingsRepo, sender whatsapp.Sender) *ReminderService {
	return &ReminderService{appointmentRepo: ar, clientRepo: cr, settingsRepo: sr, sender: sender}
}

// Start inicia a goroutine que verifica lembretes a cada minuto.
func (r *ReminderService) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				r.checkAndSend(ctx)
			}
		}
	}()
	log.Println("Serviço de lembretes iniciado")
}

func (r *ReminderService) checkAndSend(ctx context.Context) {
	settings, err := r.settingsRepo.Get(ctx)
	if err != nil {
		log.Printf("Erro ao buscar settings: %v", err)
		return
	}

	appointments, err := r.appointmentRepo.GetUpcoming(ctx, settings.ReminderMinutes)
	if err != nil {
		log.Printf("Erro ao buscar agendamentos: %v", err)
		return
	}

	for _, a := range appointments {
		msg := fmt.Sprintf("Olá %s! Lembrete: você tem um horário agendado hoje às %s. Te esperamos!", a.ClientName, a.Time)
		// ClientName vem do JOIN na query GetUpcoming
		// Precisamos do phone — buscamos pelo client_id
		phone, err := r.clientRepo.GetPhoneByID(ctx, a.ClientID)
		if err != nil {
			log.Printf("Erro ao buscar telefone do cliente %d: %v", a.ClientID, err)
			continue
		}
		if err := r.sender.SendMessage(phone, msg); err != nil {
			log.Printf("Erro ao enviar lembrete para %s: %v", phone, err)
		} else {
			log.Printf("Lembrete enviado para %s", phone)
		}
	}
}
