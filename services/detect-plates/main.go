package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/openalpr/openalpr/src/bindings/go/openalpr"
	"github.com/pubnub/go/messaging"
)

// "github.com/openalpr/openalpr/src/bindings/go/openalpr"

type payloadIn struct {
	ID          string `json:"id"`
	URL         string `json:"image_url"`
	CountryCode string `json:"countrycode"`
	Bucket      string `json:"bucket"`
}

type payloadOut struct {
	ID         string      `json:"id"`
	ImageURL   string      `json:"image_url"`
	Rectangles []rectangle `json:"rectangles"`
	Bucket     string      `json:"bucket"`
	Plate      string      `json:"plate"`
}

type rectangle struct {
	StartX int `json:"startx"`
	StartY int `json:"starty"`
	EndX   int `json:"endx"`
	EndY   int `json:"endy"`
}

func main() {
	p := new(payloadIn)
	json.NewDecoder(os.Stdin).Decode(p)

	fnStart(p.Bucket, p.ID)
	defer fnFinish(p.Bucket, p.ID)
	outfile := "working.jpg"

	alpr := openalpr.NewAlpr(p.CountryCode, "", "runtime_data")
	defer alpr.Unload()

	if !alpr.IsLoaded() {
		fmt.Println("OpenALPR failed to load!")
		return
	}
	alpr.SetTopN(10)

	downloadFile(outfile, p.URL)

	imageBytes, err := ioutil.ReadFile(outfile)
	if err != nil {
		fmt.Println(err)
	}
	results, err := alpr.RecognizeByBlob(imageBytes)

	if len(results.Plates) > 0 {
		plate := results.Plates[0]
		fmt.Printf("\n\n FOUND PLATE ------> %+v", plate)

		pout := &payloadOut{
			ID:         p.ID,
			ImageURL:   p.URL,
			Rectangles: []rectangle{{StartX: plate.PlatePoints[0].X, StartY: plate.PlatePoints[0].Y, EndX: plate.PlatePoints[2].X, EndY: plate.PlatePoints[2].Y}},
			Bucket:     p.Bucket,
			Plate:      plate.BestPlate,
		}

		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(pout)

		b2 := new(bytes.Buffer)
		json.NewEncoder(b2).Encode(pout)

		postURL := os.Getenv("FUNC_SERVER_URL") + "/draw"
		res, _ := http.Post(postURL, "application/json", b)
		fmt.Println(res.Body)

		alertPostURL := os.Getenv("FUNC_SERVER_URL") + "/alert"
		resAlert, _ := http.Post(alertPostURL, "application/json", b2)
		fmt.Println(resAlert.Body)
	} else {

		fmt.Println("No Plates Found!")

	}

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
	pubKey, subKey, ran string
	pn                  *messaging.Pubnub
	cbChannel           = make(chan []byte)
	errChan             = make(chan []byte)
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
	pn.Publish(bucket, fmt.Sprintf(`{"type":"detect-plates","running":true, "id":"%s", "runner": "%s"}`, id, os.Getenv("HOSTNAME")), cbChannel, errChan)
}

func fnFinish(bucket, id string) {
	pn.Publish(bucket, fmt.Sprintf(`{"type":"detect-plates","running":false, "id":"%s", "runner": "%s"}`, id, os.Getenv("HOSTNAME")), cbChannel, errChan)
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
