package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	// Open a file for writing
	fileName := os.Args[0]
	if fileName == "" {
		fileName = time.Now().String() + "-transactions-file"
	}

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
	numRecords := 10000

	for i := 0; i < numRecords; i++ {
		id := fmt.Sprintf("%d", i+1)
		month := rand.Intn(12) + 1 // random month (1-12)
		day := rand.Intn(28) + 1   // random day (1-28)
		date := fmt.Sprintf("%02d/%d", month, day)
		transactionValue := fmt.Sprintf("%+0.3f", rand.Float64()*500-250)

		writer.Write([]string{id, date, transactionValue})
	}
}
