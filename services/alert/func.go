package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/pubnub/go/messaging"
)

type payloadIn struct {
	ImageURL string `json:"image_url"`
	Plate    string `json:"plate"`
}

func main() {
	p := new(payloadIn)
	json.NewDecoder(os.Stdin).Decode(p)
	fnStart(os.Getenv("STORAGE_BUCKET"), p.Plate)
	defer fnFinish(os.Getenv("STORAGE_BUCKET"), p.Plate)

	outfile := "working.jpg"

	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONF_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONF_SECRET"))
	api := anaconda.NewTwitterApi(os.Getenv("TWITTER_TOKEN_KEY"), os.Getenv("TWITTER_TOKEN_SECRET"))

	timeStr := string(time.Now().Format(time.RFC3339))

	downloadFile(outfile, p.ImageURL)
	image := imgToBase64(outfile)

	media, err := api.UploadMedia(image)
	if err != nil {
		panic(err)
	}

	v := url.Values{}
	v.Set("media_ids", media.MediaIDString)

	api.PostTweet("VistaGuard Alert: Watch for license plate "+p.Plate+" [Detected "+timeStr+"]", v)
	//fmt.Println(tweet)
}

func imgToBase64(imgFile string) string {
	img, err := os.Open(imgFile)
	if err != nil {
		panic(err)
	}
	defer img.Close()

	fInfo, _ := img.Stat()
	var size int64 = fInfo.Size()
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

var (
	pubKey, subKey string
	pn             *messaging.Pubnub
	cbChannel      = make(chan []byte)
	errChan        = make(chan []byte)
)

func fnStart(bucket, id string) {
	pubKey = os.Getenv("PUBNUB_PUBLISH_KEY")
	subKey = os.Getenv("PUBNUB_SUBSCRIBE_KEY")

	pn = messaging.NewPubnub(pubKey, subKey, "", "", false, "", nil)
	go func() {
		for {
			select {
			case msg := <-cbChannel:
				fmt.Println(time.Now().Second(), ": ", string(msg))
			case msg := <-errChan:
				fmt.Println(string(msg))
			default:
			}
		}
	}()
	pn.Publish(bucket, fmt.Sprintf(`{"type":"alert","running":true, "id":"%s", "runner": "%s"}`, id, os.Getenv("HOSTNAME")), cbChannel, errChan)
}

func fnFinish(bucket, id string) {
	pn.Publish(bucket, fmt.Sprintf(`{"type":"alert","running":false, "id":"%s", "runner": "%s"}`, id, os.Getenv("HOSTNAME")), cbChannel, errChan)
	time.Sleep(time.Second * 2)
}

func init() {
	if os.Getenv("HOSTNAME") == "" {
		h, err := os.Hostname()
		if err == nil {
			os.Setenv("HOSTNAME", h)
		}
	}
}
