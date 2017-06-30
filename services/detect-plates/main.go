package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/openalpr/openalpr/src/bindings/go/openalpr"
)

func main() {
	alpr := openalpr.NewAlpr("us", "", "runtime_data")
	defer alpr.Unload()

	if !alpr.IsLoaded() {
		fmt.Println("OpenAlpr failed to load!")
		return
	}
	alpr.SetTopN(20)

	fmt.Println(alpr.IsLoaded())
	fmt.Println(openalpr.GetVersion())

	url := "http://wallpaper.pickywallpapers.com/1920x1080/red-alfa-romeo-4c-us-spec-in-the-city-back-view.jpg"
	downloadFile("sample.jpg", url)

	resultFromFilePath, err := alpr.RecognizeByFilePath("sample.jpg")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", resultFromFilePath)
	fmt.Printf("\n\n\n")

	imageBytes, err := ioutil.ReadFile("sample.jpg")
	if err != nil {
		fmt.Println(err)
	}
	resultFromBlob, err := alpr.RecognizeByBlob(imageBytes)
	fmt.Printf("%+v\n", resultFromBlob)
}

//stackoverflow.com/questions/33845770/how-do-i-download-a-file-with-a-http-request-in-go-language
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
