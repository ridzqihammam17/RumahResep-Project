package location

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	_ "sync"
	_ "sync/atomic"

	"rumah_resep/config"
	"rumah_resep/models"

	"github.com/sirupsen/logrus"
)

type Address struct {
	Street           string
	Number           int
	Neighborhood     string
	District         string
	City             string
	County           string
	State            string
	Country          string
	PostalCode       string
	FormattedAddress string
	Types            string
}

// Location structure used in the Geocoding and GeocodingReverse functions
type Location struct {
	Latitude  float64
	Longitude float64
}

var re = regexp.MustCompile(`^[-+]?([1-8]?\d(\.\d+)?|90(\.0+)?),\s*[-+]?(180(\.0+)?|((1[0-7]\d)|([1-9]?\d))(\.\d+)?)$`)

func validateLatLong(latitude, longitude string) bool {
	return re.MatchString(latitude + "," + longitude)
}

func (address *Address) FormatAddress() string {

	// Creats a slice with all content from the Address struct
	var content []string
	if address.Number > 0 {
		content = append(content, strconv.Itoa(address.Number))
	}
	content = append(content, address.Street)
	content = append(content, address.Neighborhood)
	content = append(content, address.District)
	content = append(content, address.PostalCode)
	content = append(content, address.City)
	content = append(content, address.County)
	content = append(content, address.State)
	content = append(content, address.Country)

	var formattedAddress string

	// For each value in the content slice check if it is valid
	// and add to the formattedAddress string
	for _, value := range content {
		if value != "" {
			if formattedAddress != "" {
				formattedAddress += ", "
			}
			formattedAddress += value
		}
	}
	return formattedAddress
}

// httpRequest function send the HTTP request, decode the JSON
// and return a Results structure
func httpRequest(url string) (models.Results, error) {

	var results models.Results

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return results, err
	}

	// For control over HTTP client headers, redirect policy, and other settings, create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		return results, err
	}

	// Callers should close resp.Body when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Use json.Decode for reading streams of JSON data
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return results, err
	}

	// The "OK" status indicates that no error has occurred, it means
	// the address was analyzed and at least one geographic code was returned
	if strings.ToUpper(results.Status) != "OK" {
		// If the status is not "OK" check what status was returned
		switch strings.ToUpper(results.Status) {
		case "ZERO_RESULTS":
			err = errors.New("No results found.")
			break
		case "OVER_QUERY_LIMIT":
			err = errors.New("Over your quota.")
			break
		case "REQUEST_DENIED":
			err = errors.New("Your request was denied.")
			break
		case "INVALID_REQUEST":
			err = errors.New("Probably the query is missing.")
			break
		case "UNKNOWN_ERROR":
			err = errors.New("Server error. Please, try again.")
			break
		default:
			break
		}
	}
	return results, err
}

// Geocoding function is used to convert an Address structure
// to a Location structure (latitude and longitude)
func Geocoding(address string) (Location, error) {
	var location Location
	// Convert whitespaces to +

	formattedAddress := strings.Replace(string(address), " ", "+", -1)

	// Create the URL based on the formated address
	url := config.ThirdParty.GoogleMapsGeoCodeAPIUrl + "address=" + formattedAddress

	if config.ThirdParty.GoogleMapsAPIKey != "" {
		url += "&key=" + config.ThirdParty.GoogleMapsAPIKey
	}

	// Send the HTTP request and get the results
	results, err := httpRequest(url)
	if err != nil {
		log.Println(err)
		return location, err
	}

	// Get the results (latitude and longitude)
	location.Latitude = results.Results[0].Geometry.Location.Lat
	location.Longitude = results.Results[0].Geometry.Location.Lng

	return location, nil
}

func calculateDistance(start, end []string) (result int, err error) {
	url := fmt.Sprintf(config.ThirdParty.GoogleMapsAPIUrl, start[0], start[1], end[0], end[1], config.ThirdParty.GoogleMapsAPIKey)
	defer logrus.Info(url, result)
	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	var obj map[string]interface{}
	err = json.Unmarshal(contents, &obj)
	if err != nil {
		return 0, err
	}
	rows, ok := obj["rows"].([]interface{})
	if ok == false {
		return 0, errors.New("CALCULATE_FAILED")
	}
	if len(rows) == 0 {
		return 0, errors.New("CALCULATE_FAILED")
	}
	row, ok := rows[0].(map[string]interface{})
	if ok == false {
		return 0, errors.New("CALCULATE_FAILED")
	}
	elements, ok := row["elements"].([]interface{})
	if ok == false {
		return 0, errors.New("CALCULATE_FAILED")
	}
	if len(elements) == 0 {
		return 0, errors.New("CALCULATE_FAILED")
	}
	element, ok := elements[0].(map[string]interface{})
	if ok == false {
		return 0, errors.New("CALCULATE_FAILED")
	}
	status, ok := element["status"].(string)
	if ok == false {
		return 0, errors.New("CALCULATE_FAILED")
	}
	if status == "ZERO_RESULTS" {
		return 0, nil
	}
	if status != "OK" {
		return 0, errors.New(status)
	}
	distance, ok := element["distance"].(map[string]interface{})
	if ok == false {
		return 0, errors.New("CALCULATE_FAILED")
	}
	value, ok := distance["value"].(float64)
	if ok == false {
		return 0, errors.New("CALCULATE_FAILED")
	}

	return (int)(value), nil
}
func calculateMultipleDistance(location [][]string) (result int, err error) {
  var total_distance = 0
  for i := 0; i < len(location)-1; i++{
        url := fmt.Sprintf(config.ThirdParty.GoogleMapsAPIUrl, location[i][0], location[i][1], location[i+1][0], location[i+1][1], config.ThirdParty.GoogleMapsAPIKey)
        defer logrus.Info(url, result)
	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	var obj map[string]interface{}
	err = json.Unmarshal(contents, &obj)
	if err != nil {
		return 0, err
	}
	rows, ok := obj["rows"].([]interface{})
	if ok == false {
		return 0, errors.New("CALCULATE_FAILED")
	}
	if len(rows) == 0 {
		return 0, errors.New("CALCULATE_FAILED")
	}
	row, ok := rows[0].(map[string]interface{})
	if ok == false {
		return 0, errors.New("CALCULATE_FAILED")
	}
	elements, ok := row["elements"].([]interface{})
	if ok == false {
		return 0, errors.New("CALCULATE_FAILED")
	}
	if len(elements) == 0 {
		return 0, errors.New("CALCULATE_FAILED")
	}
	element, ok := elements[0].(map[string]interface{})
	if ok == false {
		return 0, errors.New("CALCULATE_FAILED")
	}
	status, ok := element["status"].(string)
	if ok == false {
		return 0, errors.New("CALCULATE_FAILED")
	}
	if status == "ZERO_RESULTS" {
		return 0, nil
	}
	if status != "OK" {
		return 0, errors.New(status)
	}
	distance, ok := element["distance"].(map[string]interface{})
	if ok == false {
		return 0, errors.New("CALCULATE_FAILED")
	}
	value, ok := distance["value"].(float64)
	if ok == false {
		return 0, errors.New("CALCULATE_FAILED")
	}
  total_distance += (int)(value)
}
	return (int)(total_distance), nil
}
