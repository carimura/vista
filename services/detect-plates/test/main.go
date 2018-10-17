package main

import (
	"../misc"
	"C"
	"encoding/json"
	"log"
	"os"
)

func main() {
	f, err := os.Open("payload.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	var p misc.PayloadIn
	if err := json.NewDecoder(f).Decode(&p); err != nil {
		log.Fatal(err.Error())
	}
	results, err := misc.SetupALRPResults(&p)
	if err != nil {
		log.Fatal(err.Error())
	}

	if len(results.Plates) > 0 {
		pout, err := misc.ProcessALRPResulsts(results, &p)
		if err != nil {
			log.Fatal(err.Error())
		}
		misc.SaveResults(os.Stderr, pout, true)
	}
}
