package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
)

func getJson(url string, user string, pass string, target interface{}) error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("preparing request, %v", err)
	}
	req.SetBasicAuth(user, pass)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("performing request, %v", err)
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

type UBResp struct {
	Status       bool   `json:"status"`
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Data         string `json:"data"`
}

func main() {

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":false,"error_code":1,"error_message":"Client ID '23062' not found.","data":""}`))
	}))
	defer srv.Close()

	ubResp := new(UBResp)

	if err := getJson(srv.URL, "foo", "bar", &ubResp); err != nil {
		log.Fatal(err)
	}

	fmt.Println(ubResp.ErrorCode) // always prints 0 even when response has {..., "error_code": 1}
}
