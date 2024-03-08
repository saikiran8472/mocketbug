package main

import (
	"encoding/json"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func transformationCriteria(inputData map[string]interface{}) map[string]interface{} {

	resultMap := make(map[string]interface{})
	for k, v := range inputData {
		// omit fields with empty key
		if strings.TrimSpace(k) == "" {
			continue
		}
		if _, ok := v.(map[string]interface{}); !ok {
			continue
		}
		// follows transforamtion criteria for value data types
		tranforedValue := criteriaChecks(v.(map[string]interface{}))

		// omit invalid fields
		if tranforedValue != nil {
			resultMap[k] = tranforedValue
		}
	}
	return resultMap
}

func criteriaChecks(value map[string]interface{}) interface{} {

	for key := range value {
		temp := value[key]
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		value[key] = temp
	}
	for key, val := range value {

		// sanitizes trailing and leading whitespace before processing
		key = strings.TrimSpace(key)

		switch key {
		case "S":
			return transformString(val.(string))
		case "N":
			return transformNumeric(val.(string))
		case "BOOL":
			return transformBool(val.(string))
		case "NULL":
			nullValue := strings.TrimSpace(val.(string))
			if nullValue == "" {
				return nil
			}
			if _, ok := val.(string); !ok {
				continue
			}

			return transformNull(val.(string))
		case "L":
			if _, ok := val.([]interface{}); !ok {

				return nil
			}

			res := transformList(val.([]interface{}))
			if res != nil {
				return res
			}
		case "M":
			return transformMap(value)
		default:
			continue
		}
	}
	return nil
}

func transformMap(inputMap map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range inputMap {
		// sanitizes trailing and leading whitespace before processing
		if strings.TrimSpace(key) == "" {
			continue
		}
		if _, ok := value.(map[string]interface{}); !ok {
			continue
		}

		for k, val := range value.(map[string]interface{}) {

			transformedValue := criteriaChecks(val.(map[string]interface{}))

			if transformedValue != nil {
				result[k] = transformedValue
			}
		}

	}

	keys := make([]string, 0, len(result))
	for key := range result {
		keys = append(keys, key)
	}
	// lexically sort map keys
	sort.Strings(keys)

	sortedResult := make(map[string]interface{})
	for _, key := range keys {
		sortedResult[key] = result[key]
	}

	return sortedResult
}

func transformList(list []interface{}) []interface{} {

	var result []interface{}

	for _, item := range list {
		if _, ok := item.(map[string]interface{}); !ok {
			continue
		}

		transformedItem := criteriaChecks(item.(map[string]interface{}))

		// Omit fields with unsupported data types
		if transformedItem != nil {
			result = append(result, transformedItem)
		}

	}

	// omit empty lists
	if len(result) == 0 {
		return nil
	}

	return result
}

func transformNull(nullValue string) interface{} {
	// sanitizes trailing and leading whitespace before processing
	nullValue = strings.TrimSpace(nullValue)
	// ParseBool returns the boolean value represented by the string.
	// It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False. Any other value returns an error
	res, err := strconv.ParseBool(nullValue)
	// omit invalid values
	if err != nil {
		return nil
	}
	// omits false bool value
	if !res {
		return nil
	}
	return res
}

func transformBool(boolValue string) interface{} {
	// sanitizes trailing and leading whitespace before processing
	boolValue = strings.TrimSpace(boolValue)
	// omit empty bool value
	if boolValue == "" {
		return nil
	}
	// ParseBool returns the boolean value represented by the string.
	// It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False. Any other value returns an error
	result, err := strconv.ParseBool(boolValue)
	if err != nil {
		return nil
	}
	return result
}

func transformNumeric(numericValue string) interface{} {
	// sanitizes trailing and leading whitespace before processing
	numericValue = strings.TrimSpace(numericValue)
	// strip the leading zeros
	numericValue = strings.TrimLeft(numericValue, "0")
	if numericValue == "" {
		return nil
	}
	//checks for floating point data type
	if strings.Contains(numericValue, ".") {
		result, err := json.Number(numericValue).Float64()
		if err != nil {
			return nil
		}
		return result
	}

	result, err := json.Number(numericValue).Int64()
	// omit fields with invalid `Numeric` values
	if err != nil {
		return nil
	}
	return result
}

func transformString(stringValue string) interface{} {
	// sanitizes trailing and leading whitespace before processing
	stringValue = strings.TrimSpace(stringValue)
	// omit empty string values
	if stringValue == "" {
		return nil
	}
	// transform `RFC3339` formatted `Strings` to `Unix Epoch` in `Numeric` data type
	if value, err := time.Parse(time.RFC3339, stringValue); err == nil {
		return value.Unix()
	}
	return stringValue
}

func main() {
	start := time.Now()
	// Reads the json file
	inputData, err := os.ReadFile("input.json")
	if err != nil {
		log.Printf("unable to read file due to %s\n", err.Error())
	}
	var jsonData map[string]interface{}
	err = json.Unmarshal(inputData, &jsonData)
	if err != nil {
		log.Printf("unable to Unmarshall : %v", err.Error())
	}
	resOutput := transformationCriteria(jsonData)
	result, err := json.Marshal(resOutput)
	if err != nil {
		log.Printf("unable to marshal : %v", err.Error())
	}
	log.Printf("Output : %v", string(result))
	elapsed := time.Since(start)
	log.Println("Execution time : ", elapsed)
}
