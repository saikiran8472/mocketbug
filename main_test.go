package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestTransformString(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  interface{}
	}{
		{"empty string", "", nil},
		{"string with leading or trailing space", "1234  ", "1234"},
		{"string with RFC3339 to unix epoch ", "2014-07-16T20:55:46Z", int64(1405544146)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transformString(tt.input)
			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestTransformNumeric(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  interface{}
	}{
		{"floating point number", "1.50", 1.50},
		{"number with leading zero", "011  ", int64(11)},
		{"invalid number", "5215s", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transformNumeric(tt.input)
			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestTransformBool(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  interface{}
	}{
		{"invalid boolean value", "1.50", nil},
		{"bool value t with leading whitespace", " t", true},
		{"bool value 0 with trailing whitespace", "0 ", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transformBool(tt.input)
			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_transformList(t *testing.T) {
	jsonData := `[
		{
            "S": ""
          },
          {
            "N": "011"
          },
          {
            "N": "5215s"
          },
          {
            "BOOL": "f"
          },
          {
            "NULL": "0"
          }
	]`

	var results []interface{}
	json.Unmarshal([]byte(jsonData), &results)
	type args struct {
		inputList []interface{}
	}

	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "list data",
			args: args{
				inputList: results,
			},
			want: []interface{}{int64(11), false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transformList(tt.args.inputList)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_transformMap(t *testing.T) {
	jsonData := `{
    "M": {
      "bool_1": {
        "BOOL": "truthy"
      },
      "null_1": {
        "NULL ": "true"
      },
      "list_1": {
        "L": [
          {
            "S": ""
          },
          {
            "N": "011"
          },
          {
            "N": "5215s"
          },
          {
            "BOOL": "f"
          },
          {
            "NULL": "0"
          }
        ]
      }
    }
  }`
	var results map[string]interface{}
	json.Unmarshal([]byte(jsonData), &results)
	type args struct {
		inputMap map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "map data",
			args: args{
				inputMap: results,
			},
			want: map[string]interface{}{
				"list_1": []interface{}{int64(11), false},
				"null_1": true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transformMap(tt.args.inputMap)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_transformationCriteria(t *testing.T) {
	jsonData := `{
  "number_1": {
    "N": "1.50"
  },
  "string_1": {
    "S": "784498 "
  },
  "string_2": {
    "S": "2014-07-16T20:55:46Z"
  },
  "map_1": {
    "M": {
      "bool_1": {
        "BOOL": "truthy"
      },
      "null_1": {
        "NULL ": "true"
      },
      "list_1": {
        "L": [
          {
            "S": ""
          },
          {
            "N": "011"
          },
          {
            "N": "5215s"
          },
          {
            "BOOL": "f"
          },
          {
            "NULL": "0"
          }
        ]
      }
    }
  },
  "list_2": {
    "L": "noop"
  },
  "list_3": {
    "L": [
      "noop"
    ]
  },
  "": {
    "S": "noop"
  }
}`

	var results map[string]interface{}
	json.Unmarshal([]byte(jsonData), &results)
	type args struct {
		inputData map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "parsing json data",
			args: args{
				inputData: results,
			},
			want: map[string]interface{}{
				"map_1": map[string]interface{}{
					"list_1": []interface{}{int64(11), false},
					"null_1": true,
				},
				"number_1": 1.5,
				"string_1": "784498",
				"string_2": int64(1405544146),
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transformationCriteria(tt.args.inputData)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}
