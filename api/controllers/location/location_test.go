package location

import (
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
	config.InitGMapsConfig()
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

