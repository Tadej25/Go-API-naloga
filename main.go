package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/ahmetb/go-linq/v3"
	"github.com/gin-gonic/gin"
)

type CityCSV struct {
	Name         string    `json:"name"`
	Temperatures []float64 `json:"temperatures"`
	Average      float64   `json:"average"`
	Max          float64   `json:"max"`
	Min          float64   `json:"min"`
}

type City struct {
	Name               string  `json:"name"`
	AverageTemperature float64 `json:"averageTemperature"`
	Max                float64 `json:"max"`
	Min                float64 `json:"min"`
}

func (c *CityCSV) AddTemperature(temp float64) {
	c.Temperatures = append(c.Temperatures, temp)
}
func processCity(c *CityCSV) {
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

// var cityMap []City
var cityMap = make(map[string]*CityCSV)
var cityArray []City

// Main function
func main() {
	err := readCsv()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}
	router := gin.Default()
	router.GET("/cities", GetCities)
	router.GET("/city/:name", GetCityByName)
	router.GET("/AverageTemperatures", GetAverageTemperatures)
	router.POST("reload", Reload)
	router.Run("localhost:8080")
}
func GetCities(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, cityArray)
}
func GetCityByName(c *gin.Context) {
	city, err := getCityByName(c.Param("name"))
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "City not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, city)
}
func getCityByName(name string) (City, error) {
	city := From(cityArray).Where(func(c interface{}) bool {
		return c.(City).Name == name
	}).First()
	if city == nil {
		return City{}, errors.New("city not found")
	}
	result := city.(City)
	return result, nil
}
func GetAverageTemperatures(c *gin.Context) {
	// Retrieve the 'type' query parameter
	filterType := c.Query("type") // Returns the value of `type` or an empty string if not provided

	// Retrieve the 'value' query parameter
	value := c.Query("value")
	valuef, err := strconv.ParseFloat(value, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid 'value' query parameter. Please provide a valid number",
		})
		return
	}
	var above bool
	// Process the query parameters
	if filterType == "above" {
		above = true
	} else if filterType == "below" {
		above = false
	} else {
		// Handle invalid or missing type
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid or missing 'type' query parameter. Use 'above' or 'below'",
		})
		return
	}
	result, err := getAverageTemperatures(valuef, above)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Error retrieving average temperatures",
		})
		return
	}
	c.IndentedJSON(http.StatusOK, result)

}
func getAverageTemperatures(value float64, above bool) ([]City, error) {
	cities := From(cityArray).Where(func(c interface{}) bool {
		if above {
			return c.(City).AverageTemperature >= value
		}
		return c.(City).AverageTemperature <= value
	}).Results()

	if len(cities) == 0 {
		return nil, errors.New("no cities found that meet the criteria")
	}

	results := make([]City, len(cities))
	for i, city := range cities {
		results[i] = city.(City)
	}

	return results, nil
}
func Reload(c *gin.Context) {
	err := readCsv()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Data reloaded successfully"})
}

func readCsv() error {
	cityMap = make(map[string]*CityCSV)
	cityArray = nil

	path, err := readOrCreateConfig()
	if err != nil {
		return errors.New("error reading config file")
	}
	file, err := os.Open(".\\" + path)
	if err != nil {
		return errors.New("error opening file: " + path)
	}
	defer file.Close()
	fmt.Println("Reading...")
	startTime := time.Now()
	readFile(file)
	fmt.Printf("Done! Reading took %s\n", time.Since(startTime))

	cityArray = make([]City, len(cityMap))
	index := 0

	fmt.Println("Calculating...")
	startTime = time.Now()
	for _, city := range cityMap {
		processCity(city)
		cityArray[index] = City{Name: city.Name, AverageTemperature: city.Average, Max: city.Max, Min: city.Min}
		index++
		//fmt.Printf("City: %s, First Temperature: %.2f°C\n", city.Name, city.Average)
	}
	fmt.Printf("Done! Calculating took %s\n", time.Since(startTime))
	return nil
}
func readFile(file *os.File) {
	// Move the file pointer to the start of the chunk
	// Create a scanner for the current chunk
	scanner := bufio.NewScanner(file)
	// Read lines until we have processed the chunk size
	for scanner.Scan() {
		line := scanner.Text()
		processLine(line)
	}
}

func processLine(line string) {

	parts := strings.Split(line, ";")
	if len(parts) != 2 {
		return
	}
	//Pridobi ime mesta in temperaturo
	cityName := parts[0]
	temperature, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		fmt.Println("Error converting string to float:", err)
		return
	}
	//Če mesto z tem imenom že obstaja, mu dodaj novo vrednost temperaturn
	//Drugače dodaj novo mesto v mapo

	if city, exists := cityMap[cityName]; exists {
		city.AddTemperature(temperature)
	} else {
		newCity := &CityCSV{Name: cityName}
		newCity.AddTemperature(temperature)
		cityMap[cityName] = newCity
	}
}

func readOrCreateConfig() (string, error) {
	filePath := "data.config"
	defaultValue := "PATH=measaures.txt"
	// Try to open the file
	file, err := os.Open(filePath)
	if err != nil {
		// If the file doesn't exist, create it with the default value
		if os.IsNotExist(err) {
			err = os.WriteFile(filePath, []byte(defaultValue), 0644)
			if err != nil {
				return "", err
			}
			return defaultValue, nil
		}
		return "", err
	}
	defer file.Close()

	// Read the content of the file
	var path string
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		// The first line should contain the path to the measurements file
		path = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	parts := strings.Split(path, "=")

	return parts[1], nil
}
