package location

import (
	"errors"
	"testing"

	"rumah_resep/config"
)

func TestValidateLatLong(t *testing.T) {
	tt := []struct {
		latitude, longitude string
		shouldPass          bool
	}{
		{"+90.0", "-127.554334", true},
		{"45", "180", true},
		{"-90", "-180", true},
		{"-90.000", "-180.0000", true},
		{"+90", "+180", true},
		{"47.1231231", "179.99999999", true},
		{"-90.", "-180.", false},
		{"+90.1", "-100.111", false},
		{"-91", "123.456", false},
		{"045", "180", false},
		{"heap", "", false},
		{"", "", false},
	}
	for _, tc := range tt {
		result := validateLatLong(tc.latitude, tc.longitude)
		expected := tc.shouldPass
		if result != expected {
			t.Errorf("%s, %s unexpected result: got %v want %v", tc.latitude, tc.longitude, result, expected)
		}
	}
}

func TestCalculateDistance(t *testing.T) {
	config.InitConfig()
	tt := struct {
		start, end []string
		expected   int
	}{
		[]string{"40.6905615", "-73.9976592"},
		[]string{"40.6655101", "-73.89188969999998"},
		10353,
	}
	result, err := calculateDistance(tt.start, tt.end)
	if err != nil {
		t.Errorf("calculateDistance error: %v", err)
	}
	expected := tt.expected
	if result != expected {
		t.Errorf("unexpected result: got %v", result)
	}
}
func TestFormatAddress(t *testing.T) {
	var address1 string
	var address2 string
	var address3 string
	var address4 string

	address2 = "Jl. Pemuda 6, East Jakarta, Indonesia"

	address3 = "Jl. Pemuda 6, Jati, East Jakarta, DKI Jakarta, Indonesia, 13220"

	address4 = "6"

	// Table tests
	var tTests = []struct {
		address          string
		addressReceived string
	}{
		{address1,""},
		{address2,"6, Jl. Pemuda, East Jakarta, Indonesia",},
		{address3,"6"},
		{address4,"6, Jl. Pemuda, Jati, 13220, East Jakarta, DKI Jakarta, Indonesia"},
	}

	// Test with all values from the tTests
	for _, pair := range tTests {

		if pair.addressReceived != pair.address {
			t.Error("Expected:", pair.address,
				"Received:", pair.addressReceived)
		}
	}
}
func TestGeocoding(t *testing.T) {
	config.InitConfig()
	var address1 string
	var address2 string

	var location1 Location
	var location2 Location

	location1 = Location{
		Latitude:  0.0,
		Longitude: 0.0,
	}

	location2 = Location{
		Latitude:  -6.1917111,
		Longitude: 106.8897306,
	}

	address2 = "6, Jl. Pemuda, Jati, 13220, East Jakarta, DKI Jakarta, Indonesia"

	// Table tests
	var tTests = []struct {
		address  string
		location Location
		err      error
	}{
		{address1, location1, errors.New("Empty Address")},
		{address2, location2, nil},
	}

	// Test with all values from the tTests
	for _, pair := range tTests {
		location, err := Geocoding(pair.address)

		if pair.err != nil {
			if err == nil {
				t.Error("Expected:", pair.err,
					"Received: nil")
			}
		} else {
			if err != nil {
				t.Error("Expected: nil",
					"Received:", err)
			}
		}
		if location.Latitude != pair.location.Latitude {
			t.Error("Expected:", pair.location.Latitude,
				"Received:", location.Latitude)
		}
		if location.Longitude != pair.location.Longitude {
			t.Error("Expected:", pair.location.Longitude,
				"Received:", location.Longitude)
		}
	}

}
