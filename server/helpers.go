package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (s *Server) sendResponse(w http.ResponseWriter, statusCode int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) arrayToString(arr []bool) string {
	strArr := make([]string, len(arr))
	for i, v := range arr {
		if v {
			strArr[i] = "1"
		} else {
			strArr[i] = "0"
		}
	}
	return strings.Join(strArr, "")
}

func (s *Server) parseArrayString(str string) ([]bool, error) {
	if str == "" {
		return []bool{}, nil
	}

	var buffer bytes.Buffer
	for i := 0; i < len(str); i++ {
		fmt.Fprintf(&buffer, "%b", str[i])
	}
	binaryString := fmt.Sprintf("%s", buffer.Bytes())

	parts := strings.Split(binaryString, "")
	boolArray := make([]bool, len(parts))

	for i, part := range parts {
		if part == "1" {
			boolArray[i] = true
		} else if part == "0" {
			boolArray[i] = false
		} else {
			return nil, fmt.Errorf("invalid boolean value: %s", part)
		}
	}

	return boolArray, nil
}
