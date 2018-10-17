package misc

import (
	"bytes"
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
	Rectangles []rectangle `json:"rectangles"`
	Plate      string      `json:"plate"`
}

type rectangle struct {
	StartX int `json:"startx"`
	StartY int `json:"starty"`
	EndX   int `json:"endx"`
	EndY   int `json:"endy"`
}

var nextFunc = "draw"
var fnAPIURL = os.Getenv("FUNC_SERVER_URL")

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
		Rectangles: []rectangle{
			{StartX: plate.PlatePoints[0].X, StartY: plate.PlatePoints[0].Y,
				EndX: plate.PlatePoints[2].X, EndY: plate.PlatePoints[2].Y}},

		Plate: plate.BestPlate,
	}

	return pout, nil
}

func SaveResults(out io.Writer, pout *PayloadOut, noChain bool) error {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(pout)
	if err != nil {
		return err
	}

	if noChain {
		out.Write(b.Bytes())
		return nil
	} else {
		postURL := fmt.Sprintf("%s/%s", fnAPIURL, nextFunc)
		os.Stderr.WriteString(
			fmt.Sprintf("Sending %s to %s",
				string(b.Bytes()), postURL))
		res, err := http.Post(postURL, "application/json", b)
		if err != nil {
			return err
		}
		io.Copy(os.Stderr, res.Body)
		defer res.Body.Close()
	}
	return nil
}
