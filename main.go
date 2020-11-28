package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/google/uuid"
)

var WRAPPED_KEYS = map[string]bool{
	"$oid":           true,
	"$date":          true,
	"$numberDouble":  true,
	"$numberDecimal": true,
	"$numberInt":     true,
	"$numberLong":    true,
	"$binary":        true,
}

func unwrapNumber(key string, value interface{}) interface{} {
	strValue, ok := value.(string)
	if !ok {
		return value
	}

	if key == "$numberDecimal" || key == "$numberDouble" {
		if floatValue, err := strconv.ParseFloat(strValue, 64); err == nil {
			return floatValue
		}
	}
	if key == "$numberInt" || key == "$numberLong" {
		if intValue, err := strconv.ParseInt(strValue, 10, 64); err == nil {
			return intValue
		}
	}
	return value
}

func decodeBinary(value interface{}) interface{} {
	m, ok := value.(map[string]interface{})
	if !ok {
		return value
	}

	mapValue := m["base64"]
	base64Value, ok := mapValue.(string)
	if !ok {
		return value
	}

	binaryValue, err := base64.StdEncoding.DecodeString(base64Value)
	if err != nil {
		return value
	}

	subTypeValue := m["subType"]
	stringSubType, ok := subTypeValue.(string)
	if !ok {
		return value
	}

	if stringSubType == "04" {
		if uuidValue, err := uuid.FromBytes(binaryValue); err == nil {
			return uuidValue.String()
		}
	}

	return value

}

func convertValue(value interface{}, unwrapNumbers bool) interface{} {
	if m, ok := value.(map[string]interface{}); ok {
		for mapKey, mapVal := range m {
			if mapKey == "$binary" {
				return decodeBinary(mapVal)
			}

			if _, exist := WRAPPED_KEYS[mapKey]; exist {
				if unwrapNumbers {
					return unwrapNumber(mapKey, mapVal)
				}

				return mapVal
			}
			return convertMap(m, unwrapNumbers)
		}
		return convertMap(m, unwrapNumbers)
	}

	if s, ok := value.([]interface{}); ok {
		return convertSlise(s, unwrapNumbers)
	}

	return value
}

func convertMap(m map[string]interface{}, unwrapNumbers bool) map[string]interface{} {
	for key, val := range m {
		m[key] = convertValue(val, unwrapNumbers)

	}
	return m
}

func convertSlise(s []interface{}, unwrapNumbers bool) []interface{} {
	for idx, val := range s {
		s[idx] = convertValue(val, unwrapNumbers)
	}
	return s
}

func convert(s []byte, unwrapNumbers bool) []byte {
	var jsonMap map[string]interface{}

	if err := json.Unmarshal(s, &jsonMap); err != nil {
		return s
	}
	newJsonMap := convertMap(jsonMap, unwrapNumbers)

	result, err := json.Marshal(newJsonMap)
	if err != nil {
		return s
	}
	return result
}

func main() {
	inFileName := flag.String("in", "", "mongodb dump file")
	unwrapNumbers := flag.Bool("unwrap-numbers", false, "unwrap numbers")
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
		convertedJSON := convert(scanner.Bytes(), *unwrapNumbers)
		if _, err := fmt.Fprintf(os.Stdout, "%s\n", convertedJSON); err != nil {
			log.Fatal("error write result:", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
