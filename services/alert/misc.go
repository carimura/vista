package main

import (
	"encoding/base64"
	"errors"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

func postTweet(body *payloadIn) error {

	timeStr := string(time.Now().Format(time.RFC3339))

	buf, err := downloadContent(body.ImageURL)
	if err != nil {
		return err
	}

	media, err := api.UploadMedia(
		base64.StdEncoding.EncodeToString(buf))
	if err != nil {
		return err
	}

	v := url.Values{}
	v.Set("media_ids", media.MediaIDString)

	_, err = api.PostTweet("VistaGuard Alert: "+
		"Watch for license plate "+body.Plate+" [Detected "+timeStr+"]", v)

	return err
}

func setupFromEnv() error {
	if os.Getenv("HOSTNAME") == "" {
		h, err := os.Hostname()
		if err == nil {
			os.Setenv("HOSTNAME", h)
		}
	}
	consumerKey := os.Getenv("TWITTER_CONF_KEY")
	if consumerKey == "" {
		return errors.New("TWITTER_CONF_KEY not set")
	}

	consumerSecret := os.Getenv("TWITTER_CONF_SECRET")
	if consumerSecret == "" {
		return errors.New("TWITTER_CONF_SECRET not set")
	}

	accessToken := os.Getenv("TWITTER_TOKEN_KEY")
	if accessToken == "" {
		return errors.New("TWITTER_TOKEN_KEY not set")
	}
	accessTokenSecret := os.Getenv("TWITTER_TOKEN_SECRET")
	if accessTokenSecret == "" {
		return errors.New("TWITTER_TOKEN_SECRET not set")
	}

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api = anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	return nil
}

func downloadContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

