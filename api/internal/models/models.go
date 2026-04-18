package models

import "time"

type Client struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

type Appointment struct {
	ID           int       `json:"id"`
	ClientID     int       `json:"client_id"`
	ClientName   string    `json:"client_name,omitempty"`
	ServiceID    *int      `json:"service_id,omitempty"`
	ServiceName  string    `json:"service_name,omitempty"`
	ServiceNames []string  `json:"service_names,omitempty"`
	Date         string    `json:"date"`
	Time         string    `json:"time"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type Schedule struct {
	ID        int    `json:"id"`
	DayOfWeek int    `json:"day_of_week"` // 0=domingo, 6=sábado
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Active    bool   `json:"active"`
	Slot      int    `json:"slot"` // 1 a 4 (faixas de horário)
}

type BlockedSlot struct {
	ID   int    `json:"id"`
	Date string `json:"date"`
	Time string `json:"time"`
}

type Settings struct {
	ReminderMinutes int `json:"reminder_minutes"`
}

type Service struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Duration int    `json:"duration"` // minutos
	Price    int    `json:"price"`    // centavos (ex: 3000 = R$30,00)
	Active   bool   `json:"active"`
}
