package models

type City struct {
	Name               string  `json:"name"`
	AverageTemperature float64 `json:"averageTemperature"`
	Max                float64 `json:"max"`
	Min                float64 `json:"min"`
}
