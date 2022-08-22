package utilities

import (
	"encoding/csv"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func ReadData(fileName string) ([][]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return [][]string{}, err
	}

	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}

func HashPassword(pwd string) (string, error) {
	passHash, hashErr := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	return string(passHash), hashErr
}
