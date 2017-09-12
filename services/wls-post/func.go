package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

//{"address":"NYC","imageURL":"http://google.com/","owner":"John Adams","plateNumber":"ABCDEF","state":"NY"}

type payloadIn struct {
	ImageURL string `json:"imageURL"`
	Plate    string `json:"plate"`
	State    string `json:"state"`
}

type plate struct {
	Address  string `json:"address"`
	ImageURL string `json:"imageURL"`
	Owner    string `json:"owner"`
	Plate    string `json:"plateNumber"`
	State    string `json:"state"`
}

func main() {
	p := new(payloadIn)
	json.NewDecoder(os.Stdin).Decode(p)
	fmt.Printf("\npayloadIn: %+v", p)

	plate1 := plate{
		Plate:    p.Plate,
		ImageURL: p.ImageURL,
		Address:  "500 Oracle Parkway",
		State:    "CA",
		Owner:    "Chad",
	}

	payloadOut := new(bytes.Buffer)
	json.NewEncoder(payloadOut).Encode(plate1)
	fmt.Printf("\npayloadOut: %+v", payloadOut)

	postURL := os.Getenv("WLS_SERVER_URL") + "/licenseplates/rest/plates/add"

	res, err := http.Post(postURL, "application/json", payloadOut)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
}
