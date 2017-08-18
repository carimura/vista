package main

import (
	"bytes"
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

func main() {
	fmt.Println("Starting...")

	endpoint := "minio.ngrok.io"
	accessKeyID := "DEMOACCESSKEY"
	secretKeyID := "DEMOSECRETKEY"

	minioClient, err := minio.New(endpoint, accessKeyID, secretKeyID, false)
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("ffmpeg", "-i", "plate_ss.mpg", "-vf", "fps=1", "out%3d.png")
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
	for _, file := range files {
		if !strings.Contains(file.Name(), "out") {
			continue
		}
		uploadFile(minioClient, file.Name())
	}
	fmt.Println("-------------------------------------")
}

func uploadFile(m *minio.Client, filename string) {
	fmt.Println("Opening: " + filename)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Printf("Uploading...\n")
	n, err := m.FPutObject("oracle-vista-out", filename, filename, "image/png")
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
