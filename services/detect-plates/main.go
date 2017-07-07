package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/openalpr/openalpr/src/bindings/go/openalpr"
)

type payloadIn struct {
	ID            string `json:"id"`
	URL           string `json:"image_url"`
	CountryCode   string `json:"countrycode"`
	Access        string `json:"access"`
	Secret        string `json:"secret"`
	Bucket        string `json:"bucket"`
	FuncServerURL string `json:"func_server_url"`
}

type payloadOut struct {
	ID         string      `json:"id"`
	ImageURL   string      `json:"image_url"`
	Rectangles []rectangle `json:"rectangles"`
	Access     string      `json:"access"`
	Secret     string      `json:"secret"`
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
	outfile := "working.jpg"

	fmt.Printf("PayloadIn ---> %+v\n", p)

	alpr := openalpr.NewAlpr(p.CountryCode, "", "runtime_data")
	defer alpr.Unload()

	if !alpr.IsLoaded() {
		fmt.Println("OpenALPR failed to load!")
		return
	}
	alpr.SetTopN(10)

	fmt.Println("IsLoaded: " + strconv.FormatBool(alpr.IsLoaded()))
	fmt.Println("OpenALPR Version: " + openalpr.GetVersion())

	downloadFile(outfile, p.URL)

	imageBytes, err := ioutil.ReadFile(outfile)
	if err != nil {
		fmt.Println(err)
	}
	results, err := alpr.RecognizeByBlob(imageBytes)
	fmt.Println("--- Results ---")
	fmt.Printf("%T -- %+v", results, results)

	if len(results.Plates) > 0 {
		plate := results.Plates[0]
		fmt.Printf("\n\nPLATE --> %+v", plate)

		pout := &payloadOut{
			ID:         p.ID,
			ImageURL:   p.URL,
			Rectangles: []rectangle{{StartX: plate.PlatePoints[0].X, StartY: plate.PlatePoints[0].Y, EndX: plate.PlatePoints[2].X, EndY: plate.PlatePoints[2].Y}},
			Access:     p.Access,
			Secret:     p.Secret,
			Bucket:     p.Bucket,
			Plate:      plate.BestPlate,
		}

		fmt.Printf("\n\npout ---> %+v", pout)
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(pout)

		b2 := new(bytes.Buffer)
		json.NewEncoder(b2).Encode(pout)

		postURL := p.FuncServerURL + "/draw"
		res, _ := http.Post(postURL, "application/json", b)
		fmt.Println(res.Body)

		fmt.Printf("\n\nbuffer to alert ----> %+v", b2)
		alertPostURL := p.FuncServerURL + "/alert"
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
