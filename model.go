package main

import "regexp"

type ColumnDefinition struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Required bool     `json:"required"`
	Formats  []string `json:"formats"`
}

type CSVDefinition struct {
	Columns []ColumnDefinition `json:"columns"`
}

type Rule struct {
	Pattern            *regexp.Regexp
	TransformationFunc func(string) (string, error)
}

type Validator struct {
	Rules []Rule
}
