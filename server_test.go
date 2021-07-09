package main

import (
	"net/http"
	"testing"
)

//Test the servers search function
func TestSearch(t *testing.T) {
	resp, err := http.Get("127.0.0.1:8080/entriesfn")
	if err != nil {
		t.Fatalf("Get entries error = %v", err)
	}
	defer resp.Body.Close()
}
