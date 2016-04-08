package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type APIRequest struct {
	Destination string
	Origin      string
	Width       float64
	Height      float64
	Length      float64
	Weight      float64
	// MeasureType string
}

func (a APIRequest) generateRequest() RateV4Request {
	r := RateV4Request{
		Revision: 2,
		UserID:   "048NA0008090",
		Package: []Package{
			a.generatePackage("Fart"),
		},
	}
	return r
}

func (a APIRequest) generatePackage(id string) Package {
	var packageSize string
	if a.Width > 12 || a.Height > 12 || a.Length > 12 {
		packageSize = "Large"
	} else {
		packageSize = "Regular"
	}

	p := Package{
		ID:             id,
		Service:        "All",
		Size:           packageSize,
		ZipOrigination: a.Origin,
		ZipDestination: a.Destination,
		Width:          a.Width,
		Height:         a.Height,
		Length:         a.Length,
		Ounces:         a.Weight,
		Machinable:     true,
	}

	if packageSize == "Large" {
		p.Container = "Rectangular"
	}

	return p
}

type APIResponse struct {
	Rates []Postage     `json:",omitempty"`
	Error *PackageError `json:",omitempty"`
}

func checkForError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	// check that everything you need is here right here

	destinationZip := queryParams["Destination"][0]
	originZip := queryParams["Origin"][0]
	width, err := strconv.ParseFloat(queryParams["Width"][0], 64)
	checkForError(err)
	height, err := strconv.ParseFloat(queryParams["Height"][0], 64)
	checkForError(err)
	length, err := strconv.ParseFloat(queryParams["Length"][0], 64)
	checkForError(err)
	weight, err := strconv.ParseFloat(queryParams["Weight"][0], 64)
	checkForError(err)

	apiRequest := APIRequest{
		Destination: destinationZip,
		Origin:      originZip,
		Width:       width,
		Height:      height,
		Length:      length,
		Weight:      weight,
	}

	rateRequest := apiRequest.generateRequest()
	rateResponse, err := rateRequest.requestRate()
	checkForError(err)

	// apiResponse := APIResponse{
	// 	Rates: rateResponse.Package[0].Postage,
	// 	Error: rateResponse.Package[0].Error,
	// }
	var jsonResponse []byte

	if rateResponse.Package[0].Error != nil {
		failedAPIResponse := map[string]string{
			"Error": rateResponse.Package[0].Error.Description,
		}
		jsonResponse, err = json.MarshalIndent(failedAPIResponse, "", "    ")
		checkForError(err)
	} else {
		successfulAPIResponse := map[string]float32{}
		for _, v := range rateResponse.Package[0].Postage {
			successfulAPIResponse[v.MailService] = v.Rate
		}

		jsonResponse, err = json.MarshalIndent(successfulAPIResponse, "", "    ")
	}

	w.Write(jsonResponse)
}

func main() {
	http.HandleFunc("/api/rate", handleRequest)
	http.ListenAndServe(":3000", nil)
}
