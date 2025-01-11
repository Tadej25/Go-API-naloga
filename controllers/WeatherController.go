package controllers

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	linq "github.com/ahmetb/go-linq/v3"
	"github.com/gin-gonic/gin"

	_ "exmaple/Go-API-naloga/docs"
	"exmaple/Go-API-naloga/models"
)

var cityMap = make(map[string]*models.CityCSV)
var cityArray []models.City

// @Summary Get all cities
// @Description Gets the city names average temperatures max and min temperatures for all cities
// @Produce  json
// @Success 200 {array} models.City[]
// @Failure 404 {object} error
// @Router /cities [get]
func GetCities(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, cityArray)
}

// @Summary Get city by name
// @Description Gets the average temperatures max and min temperatures for a specific city
// @Produce  json
// @Param name path string true "City name"
// @Success 200 {object} models.City
// @Failure 404 {object} error
// @Router /city/{name} [get]
func GetCityByName(c *gin.Context) {
	city, err := getCityByName(c.Param("name"))
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "City not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, city)
}

func getCityByName(name string) (models.City, error) {
	city := linq.From(cityArray).Where(func(c interface{}) bool {
		return c.(models.City).Name == name
	}).First()
	if city == nil {
		return models.City{}, errors.New("city not found")
	}
	result := city.(models.City)
	return result, nil
}

// @Summary Gets average temperatures for cities above or below a given value
// @Description Gets the average temperatures for cities above or below a given value
// @Produce  json
// @Param type query string true "Type of filter (above|below)" Enums(above, below)
// @Param value query float64 true "Value for filter"
// @Success 200 {array} models.City[]
// @Failure 400 {object} error
// @Router /AverageTemperatures [get]
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
func getAverageTemperatures(value float64, above bool) ([]models.City, error) {
	cities := linq.From(cityArray).Where(func(c interface{}) bool {
		if above {
			return c.(models.City).AverageTemperature >= value
		}
		return c.(models.City).AverageTemperature <= value
	}).Results()

	if len(cities) == 0 {
		return nil, errors.New("no cities found that meet the criteria")
	}

	results := make([]models.City, len(cities))
	for i, city := range cities {
		results[i] = city.(models.City)
	}

	return results, nil
}

// @Summary Reload data from CSV file
// @Description Reloads data from the CSV file
// @Produce  json
// @Security basicAuth
// @Success 200 {string} string "Data reloaded successfully"
// @Failure 401 {object} error "Unauthorized"
// @Failure 500 {object} error "Internal Server Error"
// @Router /reload [post]
func Reload(c *gin.Context) {
	err := ReadCsv()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Data reloaded successfully"})
}
func ReadCsv() error {
	cityMap = make(map[string]*models.CityCSV)
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

	cityArray = make([]models.City, len(cityMap))
	index := 0

	fmt.Println("Calculating...")
	startTime = time.Now()
	for _, city := range cityMap {
		city.ProcessCity()
		cityArray[index] = models.City{Name: city.Name, AverageTemperature: city.Average, Max: city.Max, Min: city.Min}
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
		newCity := &models.CityCSV{Name: cityName}
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
