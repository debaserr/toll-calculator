package tolls_test

import (
	tolls "github.com/debaserr/toll-calculator"
	"testing"
	"time"
)

func TestGetTollFee(t *testing.T) {
	tc, err := tolls.NewCalculatorFromFile("test_config.json")
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		vehicle     tolls.Vehicle
		dates       []time.Time
		expectedFee int
	}{
		"car passing at 11:30": {
			vehicle:     tolls.Vehicle{VehicleType: tolls.CAR},
			dates:       []time.Time{time.Date(2019, 10, 1, 11, 30, 0, 0, time.UTC)},
			expectedFee: 8,
		},
		"car passing at 12:00": {
			vehicle:     tolls.Vehicle{VehicleType: tolls.CAR},
			dates:       []time.Time{time.Date(2019, 10, 1, 12, 0, 0, 0, time.UTC)},
			expectedFee: 0,
		},
		"car passing at 5:30, 6:05, 7:04, 7:06, 8:06, 11:00, 11:30": {
			vehicle: tolls.Vehicle{VehicleType: tolls.CAR},
			dates: []time.Time{
				time.Date(2019, 10, 1, 5, 30, 0, 0, time.UTC),
				time.Date(2019, 10, 1, 6, 5, 0, 0, time.UTC),
				time.Date(2019, 10, 1, 7, 4, 0, 0, time.UTC),
				time.Date(2019, 10, 1, 7, 6, 0, 0, time.UTC),
				time.Date(2019, 10, 1, 8, 6, 0, 0, time.UTC),
				time.Date(2019, 10, 1, 11, 0, 0, 0, time.UTC),
				time.Date(2019, 10, 1, 11, 30, 0, 0, time.UTC),
			},
			expectedFee: 44,
		},
		"car hitting max fee": {
			vehicle: tolls.Vehicle{VehicleType: tolls.CAR},
			dates: []time.Time{
				time.Date(2019, 10, 1, 6, 30, 0, 0, time.UTC),
				time.Date(2019, 10, 1, 7, 31, 0, 0, time.UTC),
				time.Date(2019, 10, 1, 8, 32, 0, 0, time.UTC),
				time.Date(2019, 10, 1, 9, 33, 0, 0, time.UTC),
				time.Date(2019, 10, 1, 10, 34, 0, 0, time.UTC),
				time.Date(2019, 10, 1, 11, 35, 0, 0, time.UTC),
			},
			expectedFee: 60,
		},
		"motorbike passing at 6:00": {
			vehicle:     tolls.Vehicle{VehicleType: tolls.MOTORBIKE},
			dates:       []time.Time{time.Date(2019, 10, 1, 6, 0, 0, 0, time.UTC)},
			expectedFee: 0,
		},
		"car passing on May first": {
			vehicle:     tolls.Vehicle{VehicleType: tolls.CAR},
			dates:       []time.Time{time.Date(2019, 5, 1, 6, 0, 0, 0, time.UTC)},
			expectedFee: 0,
		},
		"car passing on Good Friday": {
			vehicle:     tolls.Vehicle{VehicleType: tolls.CAR},
			dates:       []time.Time{time.Date(2019, 4, 19, 6, 0, 0, 0, time.UTC)},
			expectedFee: 0,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actualFee := tc.GetTollFee(test.vehicle, test.dates)
			if actualFee != test.expectedFee {
				t.Logf("\nexpected: %d\nactual: %d\n", test.expectedFee, actualFee)
				t.Fail()
			}
		})
	}
}
