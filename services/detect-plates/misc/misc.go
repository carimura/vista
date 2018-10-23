package misc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/openalpr/openalpr/src/bindings/go/openalpr"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type PayloadIn struct {
	ID          string `json:"id"`
	URL         string `json:"image_url"`
	CountryCode string `json:"countrycode"`
}

type PayloadOut struct {
	GotPlate   bool        `json:"got_plate"`
	ID         string      `json:"id"`
	ImageURL   string      `json:"image_url"`
	Rectangles []Rectangle `json:"rectangles"`
	Plate      string      `json:"plate"`
}

type Rectangle struct {
	StartX int `json:"startx"`
	StartY int `json:"starty"`
	EndX   int `json:"endx"`
	EndY   int `json:"endy"`
}

func SetupALRPResults(p *PayloadIn) (*openalpr.AlprResults, error) {
	alpr := openalpr.NewAlpr(p.CountryCode, "", "runtime_data")
	defer alpr.Unload()

	if !alpr.IsLoaded() {
		return nil, errors.New("OpenALPR failed to load!")
	}

	alpr.SetTopN(10)

	os.Stderr.WriteString(fmt.Sprintf("Checking Plate URL ---> " + p.URL))

	b, err := downloadContent(p.URL)
	if err != nil {
		os.Stderr.WriteString(
			fmt.Sprintf("Failed to download file %s: %s",
				p.URL, err))
		return nil, err
	}
	res, err := alpr.RecognizeByBlob(b)
	return &res, err
}

func downloadContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func ProcessALRPResulsts(results *openalpr.AlprResults, p *PayloadIn) (*PayloadOut, error) {
	plate := results.Plates[0]
	os.Stderr.WriteString(
		fmt.Sprintf("\n\n FOUND PLATE ------>> %+v", plate))

	pout := &PayloadOut{
		GotPlate: true,
		ID:       p.ID,
		ImageURL: p.URL,
		Rectangles: []Rectangle{
			{StartX: plate.PlatePoints[0].X, StartY: plate.PlatePoints[0].Y,
				EndX: plate.PlatePoints[2].X, EndY: plate.PlatePoints[2].Y}},

		Plate: plate.BestPlate,
	}

	return pout, nil
}

func SaveResults(out io.Writer, pout *PayloadOut) error {
	err := json.NewEncoder(out).Encode(pout)
	if err != nil {
		return err
	}
	json.NewEncoder(os.Stderr).Encode(pout)
	return nil
}
