package csv

// Column
type Column struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

// Schema
type Schema struct {
	Columns    []Column `json:"columns"`
	SkipHeader bool     `json:"skip_header"`
}
