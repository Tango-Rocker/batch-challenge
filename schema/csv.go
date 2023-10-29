package schema

// Column
type Column struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Required bool     `json:"required"`
	Formats  []string `json:"formats"`
}

// CSV
type CSV struct {
	Columns    []Column `json:"columns"`
	SkipHeader bool     `json:"skip_header"`
}
