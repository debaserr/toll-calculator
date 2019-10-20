package tolls

import (
	"encoding/json"
	"github.com/teambition/rrule-go"
	"io/ioutil"
	"time"
)

type VehicleType int

const (
	CAR VehicleType = iota + 1
	MOTORBIKE
	TRACTOR
	EMERGENCY
	DIPLOMAT
	FOREIGN
	MILITARY
)

type Vehicle struct {
	VehicleType VehicleType
}

// A Tariff represents a time rule for when a certain toll fee should be applied.
// The rruleSet specifies all occurrences of the rule and duration how long they last.
type Tariff struct {
	rruleSet *rrule.Set
	duration time.Duration
	fee      int
}

// A Calculator contains tariff and vehicle settings for calculating toll fees.
// The tariffs are ordered by the priority in which they should apply.
type Calculator struct {
	Tariffs         []Tariff
	VehicleTypeFree map[VehicleType]bool
}

// NewCalculatorFromFile creates a new Calculator from a json file.
func NewCalculatorFromFile(filename string) (*Calculator, error) {
	tc := &Calculator{}
	var cfg struct {
		VehicleTypeFree map[VehicleType]bool `json:"vehicle_type_free"`
		Tariffs         []struct {
			RruleSet string `json:"rrule"`
			Duration int64  `json:"duration"`
			Fee      int    `json:"fee"`
		}
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	tc.VehicleTypeFree = cfg.VehicleTypeFree

	for _, tariff := range cfg.Tariffs {
		set, err := rrule.StrToRRuleSet(tariff.RruleSet)
		if err != nil {
			return nil, err
		}
		tc.Tariffs = append(tc.Tariffs, Tariff{
			rruleSet: set,
			fee:      tariff.Fee,
			duration: time.Duration(tariff.Duration * int64(time.Minute)),
		})
	}

	return tc, nil
}

// GetTollFee calculates the fee for the given vehicle passing through tolls at the given dates of one day.
func (tc *Calculator) GetTollFee(vehicle Vehicle, dates []time.Time) int {
	if free, ok := tc.VehicleTypeFree[vehicle.VehicleType]; ok && free {
		return 0
	}

	totalFee := 0
	var intervalStart *time.Time
	var intervalMax int
	for i := 0; i < len(dates) && totalFee < 60; i++ {
		fee := getFeeFromDate(tc.Tariffs, dates[i])

		if fee < 1 {
			continue
		}

		if intervalStart == nil {
			intervalStart = &dates[i]
			intervalMax = fee
			totalFee += fee
			continue
		}

		if dates[i].Sub(*intervalStart).Minutes() <= 60 {
			if fee > intervalMax {
				totalFee += fee - intervalMax
				intervalMax = fee
			}
		} else {
			intervalStart = &dates[i]
			intervalMax = fee
			totalFee += fee
		}
	}

	if totalFee > 60 {
		return 60
	}

	return totalFee
}

// getFeeFromDate returns the fee of the first tariff matching the given dates
func getFeeFromDate(tariffs []Tariff, date time.Time) int {
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.AddDate(0, 0, 1)
	for _, tariff := range tariffs {
		occurrences := tariff.rruleSet.Between(dayStart, dayEnd, true)
		for _, occStart := range occurrences {
			occEnd := occStart.Add(tariff.duration)
			if (date.After(occStart) || date.Equal(occStart)) && date.Before(occEnd) {
				return tariff.fee
			}
		}
	}

	return 0
}
