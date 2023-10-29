package csv

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Validator struct {
	Rules []Rule
}

type Rule struct {
	Pattern            *regexp.Regexp
	TransformationFunc func(string) (string, error)
}

var FloatValidator = Validator{
	Rules: []Rule{
		{
			Pattern: regexp.MustCompile(`^[+-]?\d+\.\d+$`), // Matches signed floats
			TransformationFunc: func(input string) (string, error) {
				if _, err := strconv.ParseFloat(input, 64); err != nil {
					return "", err
				}
				return input, nil
			},
		},
	},
}

var IntegerValidator = Validator{
	Rules: []Rule{
		{
			Pattern: regexp.MustCompile(`^\d+$`),
			TransformationFunc: func(input string) (string, error) {
				return input, nil
			},
		},
	},
}

func validateAndTransform(value string, validator Validator) (string, error) {
	for _, rule := range validator.Rules {
		if rule.Pattern.MatchString(value) {
			return rule.TransformationFunc(value)
		}
	}
	return "", errors.New("no valid format found")
}

func dateTransformation(input string, format string) (string, error) {
	t, err := time.Parse(format, input)
	if err != nil {
		return "", err
	}
	return t.Format("2006-01-02"), nil
}

var DateValidator = Validator{
	Rules: []Rule{
		{
			Pattern: regexp.MustCompile(`^\d{1,2}-\d{4}$`), // Matches "month-year"
			TransformationFunc: func(input string) (string, error) {
				return dateTransformation(input, "01-2006")
			},
		},
		{
			Pattern: regexp.MustCompile(`^\d{4}/\d{1,2}/\d{1,2}$`), // Matches "year/day/month"
			TransformationFunc: func(input string) (string, error) {
				return dateTransformation(input, "2006/02/01")
			},
		},
		{
			Pattern: regexp.MustCompile(`^[A-Za-z]+ \d{4}$`), // Matches "January 2006"
			TransformationFunc: func(input string) (string, error) {
				return dateTransformation(input, "January 2006")
			},
		},
		{
			Pattern:            regexp.MustCompile(`^\d{1,2}/\d{1,2}$`), // Matches "month/day"
			TransformationFunc: monthDayTransformation,
		},
	},
}

func monthDayTransformation(input string) (string, error) {
	currentYear := time.Now().Year()

	// Handling for 1 or 2 digit day format
	splitInput := strings.Split(input, "/")
	if len(splitInput) != 2 {
		return "", errors.New("invalid format")
	}

	month := splitInput[0]
	day := splitInput[1]

	// Convert to a date string with current year
	dateString := fmt.Sprintf("%s/%s/%d", month, day, currentYear)

	t, err := time.Parse("1/2/2006", dateString)
	if err != nil {
		return "", err
	}

	return t.Format("2006-01-02"), nil
}
