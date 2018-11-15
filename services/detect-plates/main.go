package main

import (
	"./misc"
	"context"
	"encoding/json"
	"github.com/fnproject/fdk-go"
	"io"
	"net/http"
	"os"
)

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
	p := new(misc.PayloadIn)
	err := json.NewDecoder(in).Decode(p)
	if err != nil {
		return err
	}

	results, err := misc.SetupALRPResults(p)
	if err != nil {
		return err
	}

	if len(results.Plates) > 0 {
		pout, err := misc.ProcessALRPResulsts(results, p)
		if err != nil {
			return err
		}

		if err := misc.SaveResults(out, pout); err != nil {
			return err
		}
	} else {
		err := json.NewEncoder(out).Encode(&misc.PayloadOut{
			GotPlate: false,
			Rectangles: []misc.Rectangle{},
		})
		if err != nil {
			return err
		}
		json.NewEncoder(os.Stderr).Encode(&misc.PayloadOut{
			GotPlate: false,
			Rectangles: []misc.Rectangle{},
		})
	}
	return nil
}

func main() {
	fdk.Handle(fdk.HandlerFunc(withError))
}

func init() {
	if os.Getenv("HOSTNAME") == "" {
		h, err := os.Hostname()
		if err == nil {
			os.Setenv("HOSTNAME", h)
		}
	}
}
