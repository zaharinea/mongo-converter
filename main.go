package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

var WRAPPED_KEYS = map[string]bool{"$oid": true, "$date": true, "$numberDecimal": true, "$numberLong": true, "$binary": true}

func convertValue(value interface{}) interface{} {
	if m, ok := value.(map[string]interface{}); ok {
		for k, v := range m {
			if _, exist := WRAPPED_KEYS[k]; exist {
				return v
			}
		}
		return convertMap(m)
	}

	if s, ok := value.([]interface{}); ok {
		return convertSlise(s)
	}

	return value
}

func convertMap(m map[string]interface{}) map[string]interface{} {
	for key, val := range m {
		m[key] = convertValue(val)

	}
	return m
}

func convertSlise(s []interface{}) []interface{} {
	for idx, val := range s {
		s[idx] = convertValue(val)
	}
	return s
}

func convert(s []byte) []byte {
	var jsonMap map[string]interface{}

	if err := json.Unmarshal(s, &jsonMap); err != nil {
		return s
	}
	newJsonMap := convertMap(jsonMap)

	result, err := json.Marshal(newJsonMap)
	if err != nil {
		return s
	}
	return result
}

func main() {
	inFileName := flag.String("in", "", "mongodb dump file")
	flag.Parse()

	var scanner *bufio.Scanner
	if *inFileName == "" {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		rf, err := os.Open(*inFileName)
		if err != nil {
			log.Fatal(err)
		}
		defer rf.Close()
		scanner = bufio.NewScanner(rf)
	}

	for scanner.Scan() {
		convertedJSON := convert(scanner.Bytes())
		if _, err := fmt.Fprintf(os.Stdout, "%s\n", convertedJSON); err != nil {
			log.Fatal("error write result:", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
