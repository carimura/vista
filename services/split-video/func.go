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
)

type payloadIn struct {
	VideoFile     string `json:"video_file"`
	ServiceToCall string `json:"service_to_call"`
}

type payloadOut struct {
	ID          string `json:"id"`
	ImageURL    string `json:"image_url"`
	CountryCode string `json:"countrycode"`
}

func main() {
	fmt.Println("Starting...")

	minio_endpoint := os.Getenv("MINIO_SERVER_URL")
	accessKeyID := os.Getenv("S3_ACCESS_KEY")
	secretKeyID := os.Getenv("S3_SECRET_KEY")

	fmt.Println(minio_endpoint)

	minioClient, err := minio.New(strings.Replace(minio_endpoint, "http://", "", 1), accessKeyID, secretKeyID, false)
	if err != nil {
		log.Fatal(err)
	}
	p := new(payloadIn)
	json.NewDecoder(os.Stdin).Decode(p)

	fmt.Println(p.VideoFile)
	downloadFile("working.mp4", p.VideoFile)

	cmd := exec.Command("ffmpeg", "-i", "working.mp4", "-vf", "fps=1", "out%3d.png")
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
		url = minio_endpoint + "/videoimages/" + file.Name()
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
	fmt.Printf("\npout: %+v \n", pout)

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(pout)
	postURL := os.Getenv("FUNC_SERVER_URL") + "/detect-plates"
	_, err := http.Post(postURL, "application/json", b)
	if err != nil {
		log.Fatal(err)
	}
}

func uploadFile(m *minio.Client, filename string) {
	fmt.Println("Opening: " + filename)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Printf("Uploading...")
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
