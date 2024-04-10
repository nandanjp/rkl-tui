package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const baseUrl = "https://pokeapi.co/api/v2/"

func do(endpoint string, obj interface{}) error {
	req, err := http.NewRequest(http.MethodGet, baseUrl+endpoint, nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, &obj)
}
