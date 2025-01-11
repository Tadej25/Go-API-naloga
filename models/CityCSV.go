package models

type CityCSV struct {
	Name         string    `json:"name"`
	Temperatures []float64 `json:"temperatures"`
	Average      float64   `json:"average"`
	Max          float64   `json:"max"`
	Min          float64   `json:"min"`
}

func (c *CityCSV) AddTemperature(temp float64) {
	c.Temperatures = append(c.Temperatures, temp)
}
func (c *CityCSV) ProcessCity() {
	total := 0.0
	max := 0.0
	min := 0.0
	for _, t := range c.Temperatures {
		total += t
		if t > max || max == 0 {
			max = t
		}
		if t < min || min == 0 {
			min = t
		}
	}
	c.Average = total / float64(len(c.Temperatures))
	c.Max = max
	c.Min = min
}
