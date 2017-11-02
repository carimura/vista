package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/openalpr/openalpr/src/bindings/go/openalpr"
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

type payloadWLS struct {
	ImageURL string `json:"imageURL"`
	Plate    string `json:"plate"`
	State    string `json:"state"`
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

	outfile := "working.jpg"

	alpr := openalpr.NewAlpr(p.CountryCode, "", "runtime_data")
	defer alpr.Unload()

	if !alpr.IsLoaded() {
		fmt.Println("OpenALPR failed to load!")
		return
	}
	alpr.SetTopN(10)

	log.Println("Checking Plate URL ---> " + p.URL)
	downloadFile(outfile, p.URL)

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

		/*poutWLS := &payloadWLS{
			ImageURL: p.URL,
			Plate:    plate.BestPlate,
		}*/

		log.Printf("\n\npout: %+v ", pout)
		//log.Printf("\n\npoutWLS: %+v ", poutWLS)

		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(pout)

		b2 := new(bytes.Buffer)
		json.NewEncoder(b2).Encode(pout)

		//b3 := new(bytes.Buffer)
		//json.NewEncoder(b3).Encode(poutWLS)

		if !noChain {

			postURL := os.Getenv("FUNC_SERVER_URL") + "/draw"
			log.Printf("Sending %s to %s", string(b.Bytes()), postURL)
			res, _ := http.Post(postURL, "application/json", b)
			log.Println(res.Body)

			alertPostURL := os.Getenv("FUNC_SERVER_URL") + "/alert"
			resAlert, _ := http.Post(alertPostURL, "application/json", b2)
			fmt.Println(resAlert.Body)

			//WLSPostURL := os.Getenv("FUNC_SERVER_URL") + "/wls-post"
			//resWLSFunc, _ := http.Post(WLSPostURL, "application/json", b3)
			//fmt.Println(resWLSFunc.Body)
			defer res.Body.Close()
			defer resAlert.Body.Close()
			//defer resWLSFunc.Body.Close()
		} else {
			os.Stdout.Write(b.Bytes())
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
