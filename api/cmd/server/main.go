package main

import (
	"context"
	"log"
	"net/http"

	"barber-app/internal/config"
	"barber-app/internal/database"
	"barber-app/internal/handlers"
	"barber-app/internal/repository"
	"barber-app/internal/services"
	"barber-app/internal/whatsapp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	cfg := config.Load()

	pool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}
	defer pool.Close()

	if err := database.Migrate(pool); err != nil {
		log.Fatalf("Erro ao rodar migrations: %v", err)
	}

	// Repositories
	clientRepo := repository.NewClientRepo(pool)
	appointmentRepo := repository.NewAppointmentRepo(pool)
	scheduleRepo := repository.NewScheduleRepo(pool)
	blockedSlotRepo := repository.NewBlockedSlotRepo(pool)
	settingsRepo := repository.NewSettingsRepo(pool)
	serviceRepo := repository.NewServiceRepo(pool)

	// WhatsApp sender
	var sender whatsapp.Sender
	switch cfg.WhatsAppProvider {
	case "twilio":
		sender = whatsapp.NewTwilioSender("", "", "")
	default:
		sender = whatsapp.NewEvolutionSender(cfg.EvolutionAPIURL, cfg.EvolutionAPIKey, cfg.EvolutionInstance)
	}

	// Services
	slotService := services.NewSlotService(scheduleRepo, blockedSlotRepo, appointmentRepo)
	reminderService := services.NewReminderService(appointmentRepo, clientRepo, settingsRepo, sender)
	reminderService.Start(context.Background())

	// Handlers
	clientHandler := handlers.NewClientHandler(clientRepo)
	appointmentHandler := handlers.NewAppointmentHandler(appointmentRepo, slotService)
	adminHandler := handlers.NewAdminHandler(scheduleRepo, blockedSlotRepo, appointmentRepo, settingsRepo, cfg.AdminUser, cfg.AdminPassword)
	serviceHandler := handlers.NewServiceHandler(serviceRepo)

	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Rotas públicas (cliente)
	r.Post("/api/clients/login", clientHandler.Login)
	r.Post("/api/clients/register", clientHandler.Register)
	r.Get("/api/clients/check", clientHandler.CheckPhone)
	r.Get("/api/slots", appointmentHandler.GetAvailableSlots)
	r.Post("/api/appointments", appointmentHandler.Create)
	r.Get("/api/appointments/client/{clientID}", appointmentHandler.GetByClient)
	r.Put("/api/appointments/{id}/cancel", appointmentHandler.Cancel)
	r.Get("/api/services", serviceHandler.GetActive)

	// Rotas admin
	r.Route("/api/admin", func(r chi.Router) {
		r.Use(adminHandler.BasicAuth)
		r.Put("/schedules", adminHandler.UpsertSchedule)
		r.Get("/schedules", adminHandler.GetSchedules)
		r.Post("/blocked-slots", adminHandler.BlockSlot)
		r.Delete("/blocked-slots", adminHandler.UnblockSlot)
		r.Get("/booked-slots", adminHandler.GetBookedSlots)
		r.Get("/appointments", adminHandler.GetAppointmentsByDate)
		r.Delete("/appointments/{id}", adminHandler.CancelAppointment)
		r.Get("/settings", adminHandler.GetSettings)
		r.Put("/settings", adminHandler.UpdateSettings)
		r.Get("/services", serviceHandler.GetAll)
		r.Post("/services", serviceHandler.Create)
		r.Put("/services/{id}", serviceHandler.Update)
		r.Delete("/services/{id}", serviceHandler.Delete)
	})

	log.Printf("Servidor rodando na porta %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
