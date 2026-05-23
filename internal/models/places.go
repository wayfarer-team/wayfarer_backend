package models

type Place struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	ImageURL        string `json:"image_url"`
	Region          string `json:"region"`
	Category        string `json:"category"`
	ApproximateCost int    `json:"approximate_cost"`
}
