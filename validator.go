package main

import (
	"regexp"
	"strconv"
	"strings"
)

func floatTransformation(input string) (string, error) {
	processed := strings.ReplaceAll(input, ".", "")
	processed = strings.ReplaceAll(processed, ",", ".")
	_, err := strconv.ParseFloat(processed, 64)
	if err != nil {
		return "", err
	}
	return processed, nil
}

var FloatValidator = Validator{
	Rules: []Rule{
		{
			Pattern:            regexp.MustCompile(`^\d{1,3}(?:\.\d{3})*(?:,\d+)?$`),
			TransformationFunc: floatTransformation,
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
