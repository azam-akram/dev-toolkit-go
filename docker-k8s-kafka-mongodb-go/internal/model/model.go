package model

import "time"

type StudentRequest struct {
	CreatedBy string    `json:"createdBy,omitempty"`
	Students  []Student `json:"students,omitempty"`
}

type StudentResponse struct {
	Version   string    `json:"version,omitempty"`
	HostName  string    `json:"hostName,omitempty"`
	Created   time.Time `json:"created,omitempty"`
	CreatedBy string    `json:"createdBy,omitempty"`
	Students  []Student `json:"students,omitempty"`
}

type Student struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Year string `json:"year,omitempty"`
}
