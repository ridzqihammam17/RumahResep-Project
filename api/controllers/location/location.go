package location

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	_ "sync"
	_ "sync/atomic"
	"testing"

	"github.com/sirupsen/logrus"
	"rumah_resep/config"
)


var isTest bool

//Test use test mode
func Test() {
	isTest = true
}
var re = regexp.MustCompile(`^[-+]?([1-8]?\d(\.\d+)?|90(\.0+)?),\s*[-+]?(180(\.0+)?|((1[0-7]\d)|([1-9]?\d))(\.\d+)?)$`)

func validateLatLong(latitude, longitude string) bool {
	return re.MatchString(latitude + "," + longitude)
}

//CheckResponse CheckResponse for test
func CheckResponse(rr *httptest.ResponseRecorder, expectedCode int, expected string, t *testing.T) {
	if status := rr.Code; status != expectedCode {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedCode)
	}

	output := strings.Trim(rr.Body.String(), "\n")
	if expected != `ignore` && output != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", output, expected)
	}
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

