package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	flag.Parse()

	file, err := os.Open(flag.Args()[0])
	if err != nil {
		panic(err)
	}

	entries, err := ParseHARJson(file)
	if err != nil {
		panic(err)
	}

	Download(entries)
}

type HAR struct {
	Log struct {
		// Version string `json:"version"`
		// Creater struct {
		// 	Name    string `json:"name"`
		// 	Version string `json:"version"`
		// } `json:"creator"`
		Entries []Entries `json:"entries"`
	} `json:"log"`
}

type NameValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Entries struct {
	//StartDate string   `json:"startedDateTime"`
	//Time      float32  `json:"time"`
	Request Request `json:"request"`
	//Response  Response `json:"response"`
	//"cache": {},
	// Timings struct {
	// 	Blocked float32 `json:"blocked"`
	// 	DNS     int     `json:"dns"`
	// 	Connect int     `json:"connect"`
	// 	Send    int     `json:"send"`
	// 	Wait    float32 `json:"wait"`
	// 	Receive float32 `json:"receive"`
	// 	SSL     int     `json:"ssl"`
	// } `json:"timings"`
}

type Request struct {
	Method string `json:"method"`
	URL    string `json:"url"`
	// HttpVersion string   `json:"httpVersion"`
	// Headers     []NameValue `json:"headers"`
	QueryString []NameValue `json:"queryString"`
	// Cookies     []NameValue `json:"cookies"`
	// HeaderSize  int      `json:"headersSize"`
	// BodySize    int      `json:"bodySize"`
}

func (t Request) GetImageId() string {
	for _, query := range t.QueryString {
		if query.Name == "ImageID" {
			return query.Value
		}
	}
	return ""
}

type Response struct {
	Status      int         `json:"status"`
	StatusText  string      `json:"statusText"`
	HttpVersion string      `json:"httpVersion"`
	Headers     []NameValue `json:"headers"`
	Cookies     []NameValue `json:"cookies"`
	Content     struct {
		Size     int    `json:"content"`
		MimeType string `json:mimeType`
	} `json:"content"`
	RedirectURL  string `json:"redirectURL"`
	HeaderSize   int    `json:"headersSize"`
	BodySize     int    `json:"bodySize"`
	TransferSize int    `json:"_transferSize"`
}

func ParseHARJson(reader io.Reader) ([]Entries, error) {

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var har = new(HAR)
	if err := json.Unmarshal(b, har); err != nil {
		return nil, err
	}

	return har.Log.Entries, nil
}

func Download(entries []Entries) {
	client := http.Client{}

	for _, entry := range entries {

		if !strings.Contains(entry.Request.URL, "Getlowresimage") {
			continue
		}

		resp, err := client.Get(entry.Request.URL)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer resp.Body.Close()

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}

		filename := fmt.Sprintf("photo-%s.jpg", entry.Request.GetImageId())
		fmt.Println("Downloaded", filename)
		if err := ioutil.WriteFile(filename, data, 0644); err != nil {
			fmt.Println(err)
			continue
		}

		time.After(time.Second * 3)
	}

}
