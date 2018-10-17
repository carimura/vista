package main

import (
	"context"
	"encoding/json"
	"github.com/ChimeraCoder/anaconda"
	"github.com/fnproject/fdk-go"
	"io"
	"log"
	"net/http"
	"os"
)

var api *anaconda.TwitterApi

type payloadIn struct {
	ImageURL string `json:"image_url"`
	Plate    string `json:"plate"`
}

func withError(_ context.Context, in io.Reader, out io.Writer) {
	os.Stderr.WriteString("STARTING ALERT FUNC")
	var body payloadIn
	err := json.NewDecoder(in).Decode(&body)
	if err != nil {
		fdk.WriteStatus(out, http.StatusInternalServerError)
		out.Write([]byte(err.Error()))
		return
	}

	err = postTweet(&body)
	if err != nil {
		fdk.WriteStatus(out, http.StatusInternalServerError)
		out.Write([]byte(err.Error()))
		return
	}
	fdk.WriteStatus(out, http.StatusOK)
}

func init() {
	err := setupFromEnv()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	fdk.Handle(fdk.HandlerFunc(withError))
}
