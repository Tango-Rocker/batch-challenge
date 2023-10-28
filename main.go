package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const definitionJSON := `
	{
		"columns": [
			{
				"name": "account_id",
				"type": "integer",
				"required": true
			},
			{
				"name": "month/year",
				"type": "date",
				"required": true,
				"formats": ["01/2006", "January 2006"]
			},
			{
				"name": "amount",
				"type": "float",
				"required": true
			}
		]
	}`


func main() {
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}

	fmt.Println("using base path: "+executable)
	var def CSVDefinition
	err = json.Unmarshal([]byte(definitionJSON), &def)
	if err != nil {
		fmt.Println("Error:", err)
	}

	err = validateCSVWithDefinition("path_to_csv.csv", def)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("CSV is valid!")
	}
}

func validateCSVWithDefinition(s string, def CSVDefinition) error {
	return nil
}
