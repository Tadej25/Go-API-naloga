package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type City struct {
	Name         string    `json:"name"`
	Temperatures []float64 `json:"temperatures"`
	Average      float64   `json:"average"`
	Max          float64   `json:"max"`
	Min          float64   `json:"min"`
}

func (c *City) AddTemperature(temp float64) {
	c.Temperatures = append(c.Temperatures, temp)
}
func processCity(c *City) {
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
}

// var cityMap []City
var cityMap = make(map[string]*City)
var cityArray []City

// Main function
func main() {
	// readCsv()
	router := gin.Default()
	router.GET("/test", TestApiCall)
	router.Run("localhost:8080")
}
func TestApiCall(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "API is working!"})
}
func readCsv() {
	path, err := readOrCreateConfig()
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	file, err := os.Open(".\\" + path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
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
		cityArray[index] = *city
		index++
		processCity(city)
		//fmt.Printf("City: %s, First Temperature: %.2f°C\n", city.Name, city.Average)
	}
	fmt.Printf("Done! Calculating took %s\n", time.Since(startTime))
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
		newCity := &City{Name: cityName}
		newCity.AddTemperature(temperature)
		cityMap[cityName] = newCity
	}
}

func readOrCreateConfig() (string, error) {
	filePath := "data.config"
	defaultValue := "PATH=measures.txt"
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
		path = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	parts := strings.Split(path, "=")

	return parts[1], nil
}
