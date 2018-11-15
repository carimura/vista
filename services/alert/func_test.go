package main

import (
	"encoding/json"
	"os"
	"testing"
)


// this test basically ensures whether a function can post tweets with keys in env vars
func TestHandler(t *testing.T) {
	t.Run("twitter-post-test", func(t *testing.T) {
		// where you must set twitter config through env vars:
		// - TWITTER_CONF_KEY
		// - TWITTER_CONF_SECRET
		// - TWITTER_TOKEN_KEY
		// - TWITTER_TOKEN_SECRET
		err := setupFromEnv()
		if err != nil {
			t.Fatal(err.Error())
		}

		payloadFile, err := os.Open("payload.json")
		if err != nil {
			t.Fatal(err.Error())
		}
		var body payloadIn
		if err := json.NewDecoder(payloadFile).Decode(&body); err != nil {
			t.Fatal(err.Error())
		}
		if err = postTweet(&body); err != nil {
			t.Fatal(err.Error())
		}
	})
}
