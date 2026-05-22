package models

type Event struct {
	ID           int    `json:"id"`
	DayID        int    `json:"day_id"`
	Title        string `json:"title"`
	StartTime    string `json:"start_time"`
	LocationName string `json:"location_name"`
	Cost         int    `json:"cost"`
	Description  string `json:"description"`
	TripID       int    `json:"trip_id"`
}
