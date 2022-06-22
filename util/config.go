package util

import (
	"encoding/csv"
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func ReadConfig(configFilePath string) kafka.ConfigMap {
	m := make(map[string]kafka.ConfigValue)
	file, err := os.Open(configFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Bad file : %s", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "#") && len(line) != 0 {
			kv := strings.Split(line, "=")
			parameter := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			m[parameter] = value
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Failed to read from file: %s", err)
		os.Exit(1)
	}
	return m
}

func ReadCSV(csvFilePath string) [][]string {
	file, err := os.Open(csvFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Bad file : %s", err)
		os.Exit(1)
	}
	defer file.Close()
	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading csv file: %s", err)
		os.Exit(1)
	}
	return data
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist) 
}
