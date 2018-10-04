package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fnproject/fdk-go"
	"github.com/openalpr/openalpr/src/bindings/go/openalpr"
	"io"
	"io/ioutil"
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
	p := new(payloadIn)
	err := json.NewDecoder(in).Decode(p)
	if err != nil {
		return err
	}

	_, noChain := os.LookupEnv("NO_CHAIN")
	if noChain {
		os.Stderr.WriteString("running without chaining")
	}

	alpr := openalpr.NewAlpr(p.CountryCode, "", "runtime_data")
	defer alpr.Unload()

	if !alpr.IsLoaded() {
		os.Stderr.WriteString("OpenALPR failed to load!")
		return errors.New("OpenALPR failed to load!")
	}

	alpr.SetTopN(10)

	os.Stderr.WriteString(fmt.Sprintf("Checking Plate URL ---> " + p.URL))

	b, err := downloadContent(p.URL)
	if err != nil {
		os.Stderr.WriteString(
			fmt.Sprintf("Failed to download file %s: %s",
				p.URL, err))
		return err
	}

	results, err := alpr.RecognizeByBlob(b)
	if err != nil {
		return err
	}

	if len(results.Plates) > 0 {
		plate := results.Plates[0]
		os.Stderr.WriteString(
			fmt.Sprintf("\n\n FOUND PLATE ------>> %+v", plate))

		pout := &payloadOut{
			GotPlate: true,
			ID:       p.ID,
			ImageURL: p.URL,
			Rectangles: []rectangle{
				{StartX: plate.PlatePoints[0].X, StartY: plate.PlatePoints[0].Y,
					EndX: plate.PlatePoints[2].X, EndY: plate.PlatePoints[2].Y}},

			Plate: plate.BestPlate,
		}

		os.Stderr.WriteString(fmt.Sprintf("\n\npout: %+v ", pout))

		b := new(bytes.Buffer)
		err = json.NewEncoder(b).Encode(pout)
		if err != nil {
			return err
		}

		if noChain {
			os.Stdout.Write(b.Bytes())
		} else {
			postURL := os.Getenv("FUNC_SERVER_URL") + "/draw"
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
	} else {
		os.Stderr.WriteString("No Plates Found!")
		if noChain {
			err := json.NewEncoder(os.Stdout).Encode(&payloadOut{
				GotPlate: false,
			})
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func main() {
	fdk.Handle(fdk.HandlerFunc(withError))
}

func downloadContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func init() {
	if os.Getenv("HOSTNAME") == "" {
		h, err := os.Hostname()
		if err == nil {
			os.Setenv("HOSTNAME", h)
		}
	}
}
