package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {

	fileName := "data-" + time.Now().Format("2006-01-02:15:04")
	println("Generating file: " + fileName)
	file, err := os.Create(fileName + ".csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	writer.Write([]string{"id", "date", "amount"})

	// Number of records to generate
	numRecords := 100000

	for i := 0; i < numRecords; i++ {
		id := fmt.Sprintf("%d", i+1)
		month := rand.Intn(2) + 6 // random month (1-12)
		day := rand.Intn(28) + 1  // random day (1-28)
		date := fmt.Sprintf("%02d/%d", month, day)
		transactionValue := fmt.Sprintf("%+0.3f", rand.Float64()*500-250)

		writer.Write([]string{id, date, transactionValue})
	}
}
