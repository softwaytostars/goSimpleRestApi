package models

type Document struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description"`
}
