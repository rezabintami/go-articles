package request

import (
	"go-articles/modules/mountains"
	"time"
)

type Mountains struct {
	Name           string    `json:"name" validate:"required" validName:"name"`
	Description    string    `json:"description" validate:"required" validName:"description"`
	About          string    `json:"about" validate:"required" validName:"about"`
	Status         string    `json:"status" validate:"required" validName:"status"`
	Latitude       float64   `json:"latitude" validate:"required" validName:"latitude"`
	Longitude      float64   `json:"longitude" validate:"required" validName:"longitude"`
	Province       string    `json:"province" validate:"required" validName:"province"`
	Country        string    `json:"country" validate:"required" validName:"country"`
	Type           string    `json:"type" validate:"required" validName:"type"`
	Height         int       `json:"height" validate:"required" validName:"height"`
	Difficult      string    `json:"difficult" validate:"required" validName:"difficult"`
	LastEruption   time.Time `json:"last_eruption" validate:"required" validName:"last_eruption"`
	TemperatureMin float64   `json:"temperature_min" validate:"required" validName:"temperature_min"`
	TemperatureMax float64   `json:"temperature_max" validate:"required" validName:"temperature_max"`
}

func (request *Mountains) ToDomain() *mountains.Domain {
	return &mountains.Domain{
		Name:           request.Name,
		Description:    request.Description,
		About:          request.About,
		Status:         request.Status,
		Latitude:       request.Latitude,
		Longitude:      request.Longitude,
		Province:       request.Province,
		Country:        request.Country,
		Type:           request.Type,
		Height:         request.Height,
		Difficult:      request.Difficult,
		LastEruption:   request.LastEruption,
		TemperatureMin: request.TemperatureMin,
		TemperatureMax: request.TemperatureMax,
	}
}
