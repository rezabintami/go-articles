package response

import (
	"go-articles/modules/mountains"
	"time"
)

type Mountains struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	About          string    `json:"about"`
	Status         string    `json:"status"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	Province       string    `json:"province"`
	Country        string    `json:"country"`
	Type           string    `json:"type"`
	Height         int       `json:"height"`
	Difficult      string    `json:"difficult"`
	LastEruption   time.Time `json:"last_eruption"`
	TemperatureMin float64   `json:"temperature_min"`
	TemperatureMax float64   `json:"temperature_max"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type MountainsFetch struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Status         string    `json:"status"`
	Province       string    `json:"province"`
	Country        string    `json:"country"`
	Type           string    `json:"type"`
	Height         int       `json:"height"`
	Difficult      string    `json:"difficult"`
	LastEruption   time.Time `json:"last_eruption"`
	TemperatureMin float64   `json:"temperature_min"`
	TemperatureMax float64   `json:"temperature_max"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type MountainsList struct {
	Mountains *[]MountainsFetch `json:"mountains"`
	Total     int               `json:"total"`
}

func FromDomain(domain mountains.Domain) Mountains {
	return Mountains{
		ID:             domain.ID,
		Name:           domain.Name,
		Description:    domain.Description,
		About:          domain.About,
		Status:         domain.Status,
		Latitude:       domain.Latitude,
		Longitude:      domain.Longitude,
		Province:       domain.Province,
		Country:        domain.Country,
		Type:           domain.Type,
		Height:         domain.Height,
		Difficult:      domain.Difficult,
		LastEruption:   domain.LastEruption,
		TemperatureMin: domain.TemperatureMin,
		TemperatureMax: domain.TemperatureMax,
		CreatedAt:      domain.CreatedAt,
		UpdatedAt:      *domain.UpdatedAt,
	}
}

func FromListDomain(domain []mountains.Domain) *[]Mountains {
	result := []Mountains{}
	for _, value := range domain {
		article := Mountains{
			ID:             value.ID,
			Name:           value.Name,
			Description:    value.Description,
			About:          value.About,
			Status:         value.Status,
			Latitude:       value.Latitude,
			Longitude:      value.Longitude,
			Province:       value.Province,
			Country:        value.Country,
			Type:           value.Type,
			Height:         value.Height,
			Difficult:      value.Difficult,
			LastEruption:   value.LastEruption,
			TemperatureMin: value.TemperatureMin,
			TemperatureMax: value.TemperatureMax,
			CreatedAt:      value.CreatedAt,
			UpdatedAt:      *value.UpdatedAt,
		}
		result = append(result, article)
	}

	return &result
}

func FetchFromListDomain(domain []mountains.Domain, count int) *MountainsList {
	mountainsList := []MountainsFetch{}
	for _, value := range domain {
		mountain := MountainsFetch{
			ID:             value.ID,
			Name:           value.Name,
			Description:    value.Description,
			Status:         value.Status,
			Province:       value.Province,
			Country:        value.Country,
			Type:           value.Type,
			Height:         value.Height,
			Difficult:      value.Difficult,
			LastEruption:   value.LastEruption,
			TemperatureMin: value.TemperatureMin,
			TemperatureMax: value.TemperatureMax,
			UpdatedAt:      *value.UpdatedAt,
		}
		mountainsList = append(mountainsList, mountain)
	}
	result := MountainsList{}
	result.Mountains = &mountainsList
	result.Total = count
	return &result
}
