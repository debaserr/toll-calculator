# Toll fee calculator
A calculator for vehicle toll fees.

## How to use
The calculator is easiest to set up with a json config file containing the tariff and vehicle settings.
```
calculator, err := tolls.NewCalculatorFromFile("config.json")
```
Check out `test_config.json` for an example of how it should look.

The tariff time rules are specified with [rrule](https://icalendar.org/iCalendar-RFC-5545/3-8-5-3-recurrence-rule.html).
See [this demo](https://jakubroztocil.github.io/rrule/) to get an idea of how it works.
  
## Running tests
To run the tests on your machine, just run:
```
GO111MODULE=on go test ./...
```
in your terminal.
