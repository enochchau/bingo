package api

import (
	"strings"
	"testing"
)

type TestJson struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestDecodeJson(t *testing.T) {
	jsonStr := "{\"name\":\"bingo\",\"age\":5}"
	body := strings.NewReader(jsonStr)

	payload, err := decodeJson[TestJson](body)
	if err != nil {
		t.Fatalf("Failed to decode message: %v", err)
	}

	if payload.Age != 5 {
		t.Fatalf(
			"Failed to properly decode json string %s, expected %d but got %d",
			jsonStr,
			5,
			payload.Age,
		)

	}
	if payload.Name != "bingo" {
		t.Fatalf(
			"Failed to properly decode json string %s, expected '%s', but got '%s'",
			jsonStr,
			"bingo",
			payload.Name,
		)
	}
}

func TestCreateJsonError(t *testing.T) {
	errMsg := "there was an error"
	j, err := CreateErrorJson(errMsg)

	if err != nil {
		t.Fatalf("Failed to create json: %v", err)
	}
	errJson := "{\"error\":\"" + errMsg + "\"}"
	if string(j) != errJson {
		t.Fatalf("Failed to create error message, expected %s but got %s", errJson, string(j))
	}
}
