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
	"os/exec"
	"strings"

	"github.com/minio/minio-go"
)

type payloadIn struct {
	ID          string `json:"id"`
	URL         string `json:"image_url"`
	CountryCode string `json:"countrycode"`
}

type payloadOut struct {
	ID          string `json:"id"`
	ImageURL    string `json:"image_url"`
	CountryCode string `json:"countrycode"`
}

func main() {
	fmt.Println("Starting...")

	endpoint := "minio.ngrok.io"
	accessKeyID := "DEMOACCESSKEY"
	secretKeyID := "DEMOSECRETKEY"

	minioClient, err := minio.New(endpoint, accessKeyID, secretKeyID, false)
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("ffmpeg", "-i", "traffic.mp4", "-vf", "fps=1", "out%3d.png")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("------ Out files after split ------")
	url := ""
	for _, file := range files {
		if !strings.Contains(file.Name(), "out") {
			continue
		}
		uploadFile(minioClient, file.Name())
		url = "http://minio.ngrok.io/videoimages/" + file.Name()
		callDetectPlates(file.Name(), url)
	}
	fmt.Println("-------------------------------------")
}

func callDetectPlates(f string, url string) {
	pout := &payloadOut{
		ID:          f,
		ImageURL:    url,
		CountryCode: "eu",
	}
	fmt.Printf("\n\npout! --> %+v ", pout)

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(pout)

	//postURL := os.Getenv("FUNC_SERVER_URL") + "/detect-plates"
	postURL := "http://fnlocal.ngrok.io/r/myapp/detect-plates"
	res, _ := http.Post(postURL, "application/json", b)
	fmt.Println(res.Body)
}

func uploadFile(m *minio.Client, filename string) {
	fmt.Println("Opening: " + filename)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Printf("Uploading...\n")
	n, err := m.FPutObject("videoimages", filename, filename, "image/png")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully uploaded file %q of size %d\n", filename, n)
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
