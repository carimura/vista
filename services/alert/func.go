package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/ChimeraCoder/anaconda"
	"github.com/fnproject/fdk-go"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

var api *anaconda.TwitterApi

type payloadIn struct {
	ImageURL string `json:"image_url"`
	Plate    string `json:"plate"`
}

func withError(ctx context.Context, in io.Reader, out io.Writer) {
	err := myHandler(ctx, in, out)
	if err != nil {
		fdk.WriteStatus(out, http.StatusInternalServerError)
		out.Write([]byte(err.Error()))
		return
	}
	fdk.WriteStatus(out, http.StatusOK)
}

func myHandler(_ context.Context, in io.Reader, out io.Writer) error {
	os.Stderr.WriteString("STARTING ALERT FUNC")
	var body payloadIn
	err := json.NewDecoder(in).Decode(&body)
	if err != nil {
		return err
	}

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

func init() {

	if os.Getenv("HOSTNAME") == "" {
		h, err := os.Hostname()
		if err == nil {
			os.Setenv("HOSTNAME", h)
		}
	}

	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONF_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONF_SECRET"))
	api = anaconda.NewTwitterApi(os.Getenv("TWITTER_TOKEN_KEY"), os.Getenv("TWITTER_TOKEN_SECRET"))
}

func downloadContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func main() {
	fdk.Handle(fdk.HandlerFunc(withError))
}
