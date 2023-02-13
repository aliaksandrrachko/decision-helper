package entities

import "time"

type Decision struct {
	Id             int64     `json:"id"`
	Title          string    `json:"title"`
	NavigationUser string    `json:"navigationUser"`
	NavigationDate time.Time `json:"navigationDate"`
}
