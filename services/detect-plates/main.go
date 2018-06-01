package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/openalpr/openalpr/src/bindings/go/openalpr"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type payloadIn struct {
	ID          string `json:"id"`
	URL         string `json:"image_url"`
	CountryCode string `json:"countrycode"`
}

type payloadOut struct {
	GotPlate   bool        `json:"got_plate"`
	ID         string      `json:"id"`
	ImageURL   string      `json:"image_url"`
	Rectangles []rectangle `json:"rectangles"`
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

	_, noChain := os.LookupEnv("NO_CHAIN")
	if noChain {
		log.Println("running without chaining")
	}

	outfile := "/tmp/working.jpg"

	alpr := openalpr.NewAlpr(p.CountryCode, "", "runtime_data")
	defer alpr.Unload()

	if !alpr.IsLoaded() {
		fmt.Println("OpenALPR failed to load!")
		return
	}
	alpr.SetTopN(10)

	log.Println("Checking Plate URL ---> " + p.URL)
	err := downloadFile(outfile, p.URL)
	if err != nil {
		log.Fatalf("Failed to download file %s: %s", p.URL, err)
	}

	imageBytes, err := ioutil.ReadFile(outfile)
	if err != nil {
		fmt.Println(err)
	}
	results, err := alpr.RecognizeByBlob(imageBytes)

	if len(results.Plates) > 0 {
		plate := results.Plates[0]
		log.Printf("\n\n FOUND PLATE ------>> %+v", plate)

		pout := &payloadOut{
			GotPlate:   true,
			ID:         p.ID,
			ImageURL:   p.URL,
			Rectangles: []rectangle{{StartX: plate.PlatePoints[0].X, StartY: plate.PlatePoints[0].Y, EndX: plate.PlatePoints[2].X, EndY: plate.PlatePoints[2].Y}},
			Plate:      plate.BestPlate,
		}

		log.Printf("\n\npout: %+v ", pout)

		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(pout)

		if noChain {
			os.Stdout.Write(b.Bytes())
		} else {
			postURL := os.Getenv("FUNC_SERVER_URL") + "/draw"
			log.Printf("Sending %s to %s", string(b.Bytes()), postURL)
			res, _ := http.Post(postURL, "application/json", b)
			log.Println(res.Body)

			defer res.Body.Close()
		}
	} else {
		log.Println("No Plates Found!")
		if noChain {
			json.NewEncoder(os.Stdout).Encode(&payloadOut{
				GotPlate: false,
			})
		}

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

func init() {
	if os.Getenv("HOSTNAME") == "" {
		h, err := os.Hostname()
		if err == nil {
			os.Setenv("HOSTNAME", h)
		}
	}
}
