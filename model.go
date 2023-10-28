package main

type ColumnDefinition struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Required bool     `json:"required"`
	Formats  []string `json:"formats"`
}

type CSVDefinition struct {
	Columns []ColumnDefinition `json:"columns"`
}
