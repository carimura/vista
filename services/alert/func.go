package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/ChimeraCoder/anaconda"
	"github.com/fnproject/fdk-go"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const outfile string = "/tmp/working.jpg"

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
	}
}

func myHandler(_ context.Context, in io.Reader, out io.Writer) error {
	os.Stderr.WriteString("STARTING ALERT FUNC")
	var body payloadIn
	err := json.NewDecoder(os.Stdin).Decode(&body)
	if err != nil {
		return err
	}

	timeStr := string(time.Now().Format(time.RFC3339))

	err = downloadFile(outfile, body.ImageURL)
	if err != nil {
		return err
	}

	image := imgToBase64(outfile)

	media, err := api.UploadMedia(image)
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

func imgToBase64(imgFile string) string {
	img, err := os.Open(imgFile)
	if err != nil {
		panic(err)
	}
	defer img.Close()

	fInfo, _ := img.Stat()
	size := fInfo.Size()
	buf := make([]byte, size)
	fReader := bufio.NewReader(img)
	fReader.Read(buf)
	imgBase64Str := base64.StdEncoding.EncodeToString(buf)

	return imgBase64Str
}

func downloadFile(filepath string, url string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	fdk.Handle(fdk.HandlerFunc(withError))
}
